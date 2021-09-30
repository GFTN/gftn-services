// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
//
import { Injectable } from '@angular/core';
import { Spec, Schema, Operation, Path } from 'swagger-schema-official';
import {
  merge, forEach, has, get, set, isEmpty, toArray, isArray, split,
  isObject, cloneDeep, isUndefined, isString, mapValues, find, last, replace
} from 'lodash';
import { HttpClient } from '@angular/common/http';
import { VERSION_DETAILS, IWWApis } from '../../../shared/constants/versions.constant';
import { VersionService } from './version.service';
import { Router } from '@angular/router';
import { MarkdownService } from 'ngx-markdown';

interface INavPathsGroups {
  [groupName: string]: INavPathDef[];
}

export interface INavPathDef {
  groupName: string;
  pathName: string;
  operationType: string;
  operation: Operation;
  jumpText: string;
}

@Injectable()
export class SwaggerService {

  // set the _jsonSpec from the ApiComponent
  _jsonSpec: any;

  // used to add missing models to json spec
  private missingModels: { [definitionsName: string]: Schema } = {};

  navPaths: INavPathsGroups = {};

  versionConst: IWWApis = VERSION_DETAILS;

  apiFileName: string;

  constructor(
    private http: HttpClient,
    private versionService: VersionService,
    private router: Router,
    private markdownService: MarkdownService
  ) {
    this._jsonSpec = {} as any;
  }

  /**
   * Used by swagger.guard.ts to populated docs based on selected option
   *
   * @param {string} apiName
   * @returns
   * @memberof SwaggerService
   */
  async populateSwaggerDocs(apiName: string) {
    await this.getOpenApiJson(apiName);
    // await this.addOmittedDefinitions();
    return;
  }

  /**
   * Retrieves the OpenApi File based on domain
   *
   * @param {IWWApi} selectedApi - api info
   * @param {string} apiName - the file name of the API definition in the /assets folder
   * @returns
   * @memberof ApiComponent
   */
  async getOpenApiJson(apiName: string): Promise<Spec> {

    // set api File name to selected json filename via path param
    this.apiFileName = apiName;

    // get swagger def file over http
    const swaggerDef: Spec = await this.http.get(
      '/assets/open-api/' +
      this.versionService.current.module +
      '/' + apiName + '.json'
    ).toPromise() as any;

    // console.log(this.navPaths);

    this._jsonSpec = swaggerDef;

    // add missing models to json spec:
    // IMPORTANT: run deepRecursiveAddMissingModels() before deepRecursiveSwaggerUpdate()
    // so that nested models have general swagger updates applied to missing models too
    this.deepRecursiveAddMissingModels(this._jsonSpec.definitions);

    // console.log(this.missingModels);

    // combine missing models with existing models
    this._jsonSpec.definitions = merge(this._jsonSpec.definitions, this.missingModels);

    // set default props with descriptions converted to markdown:
    this._jsonSpec = this.deepRecursiveSwaggerUpdate(cloneDeep(swaggerDef));

    // console.log(this._jsonSpec);

    // set nav-paths
    this.groups(this._jsonSpec);

    return;

  }

  /**
   * Assign endpoints under category grouping
   *
   * @private
   * @param {Spec} swaggerDef
   * @param {SwaggerService} self
   * @returns
   * @memberof SwaggerService
   */
  private groups(swaggerDef: Spec) {

    this.navPaths = {};

    // loop each path
    forEach(swaggerDef.paths, (pathObj: Path, pathName: string) => {

      // loop each operation in each path
      forEach(pathObj, (operationObj: Operation, operationName: string) => {

        // if no group provided don't include in array of navPaths
        if (operationObj['x-group-e']) {

          // get the root url as the group name
          const group = operationObj['x-group-e'];

          // create group if doesn't exist
          if (!has(this.navPaths, group)) {
            // create group if it does not exist
            set(this.navPaths, group, []);
          }

          // create object
          this.navPaths[group].push({
            groupName: group,
            pathName: pathName,
            operationType: operationName,
            operation: operationObj,
            jumpText: this.jumpText('path', operationName, pathName)
          });

        }

      });

    });

    return;

  }

  /**
   * Adds models that are missing from /definitions/xxxx
   * as a result of being deeply nested
   *
   * @param {{}} obj
   * @returns
   * @memberof SwaggerService
   */
  public deepRecursiveAddMissingModels(obj: {}) {

    return mapValues(obj, (o) => {

      // only apply updates to objects
      // NOTE: isObject([]) equals true so need to check
      // that is also NOT and array to confirm that o is
      // truly an object
      if (isObject(o)) {

        // CONDITIONAL MODIFICATION 1:
        // check for properties OBJECT nested in another properties obj
        // (ie: a model inside a model, hence that model would not be
        // displayed in the view if looped in the ui)
        // previously referred to as adding missing models
        if (has(o, 'properties')) {
          const _jsonSpec = this._jsonSpec;
          // check if properties exist on /definitions
          forEach(o['properties'], (val: Schema, key: string) => {
            // check if is a nested properties
            if (has(val, 'properties')) {
              // has a nested properties

              if (isUndefined(get(val, 'title'))) {
                console.error('Error: model with missing title: ', { key: key, val: val });
              }

              // set title
              const title = get(val, 'title');
              // check if matching model already exists
              if (
                !isEmpty(title) &&
                isEmpty(_jsonSpec.definitions[title])
                && isEmpty(this.missingModels[title])
              ) {
                // missing so add model to /definitions
                set(this.missingModels, title, val);
              }

            }
          });

        }

        // CONDITIONAL MODIFICATION 2:
        // check for ARRAY of properties obj nested in another properties obj
        // (ie: a model inside a model, hence that model would not be
        // displayed in the view if looped in the ui)
        // previously referred to as adding missing models
        if (get(o, 'items.type') === 'object') {
          const _jsonSpec = this._jsonSpec;
          // check if is a nested properties

          if (isUndefined(get(o, 'items.title'))) {
            console.error('Error: model items with missing title: ', o['items']);
          }

          // check if matching model already exists
          const item = get(o, 'items');
          if (
            !isEmpty(get(item, 'title')) &&
            isEmpty(_jsonSpec.definitions[item.title]) &&
            isEmpty(this.missingModels[item.title])
          ) {
            // missing so add model to /definitions
            set(this.missingModels, item.title, item);
          }

        }

        // recursively apply the same (ie: convert desc to markdown)
        // operations on the object's nested structure
        this.deepRecursiveAddMissingModels(o);

      }

      return;

    });

  }

  /**
   * Updates model recursively for certain conditions
   *
   * @param {{}} obj
   * @returns
   * @memberof SwaggerService
   */
  public deepRecursiveSwaggerUpdate(obj: {}) {

    return mapValues(obj, (o) => {

      // only apply updates to objects
      // NOTE: isObject([]) equals true so need to check
      // that is also NOT and array to confirm that o is
      // truly an object
      if (isObject(o)) {

        // add swagger conditional modifications below:

        // CONDITIONAL MODIFICATION 1:
        // check if this object has a description property
        if (has(o, 'description')) {
          // update the description if is string
          if (isString(o['description'])) {
            // o['description'] = 'desc here'; // for testing, see where description is being replaced in view
            o['description'] = this.toMarkdown(o['description']); // for testing, see where description is being replaced in view
          }
        }

        if (has(o, 'schema')) {
          o['test'] = 'test';
        }

        // CONDITIONAL MODIFICATION 2:
        // add model for every schema.ref as 'object'
        if (has(o, '$ref') || has(o, 'schema.$ref')) {

          const ref = get(o, '$ref') || get(o, 'schema.$ref');

          // consistent with yaml model name
          const key = this.refName(ref) as string;
          o['model'] = this._jsonSpec.definitions[key];
        }

        // CONDITIONAL MODIFICATION 3:
        // add model for every schema.ref as 'array of objects'
        if (has(o, 'schema.items.$ref')) {
          // consistent with yaml model name
          const key = this.refName(get(o, 'schema.items.$ref')) as string;
          o['model'] = this._jsonSpec.definitions[key];
        }

        // CONDITIONAL MODIFICATION 4:
        // check if properties (specific field in the model) is required
        // required is determined by including the key of the property in
        // the required array
        if (isArray(get(o, 'required'))) {
          // required is included in the swagger spec else where other than in definition properties
          // but required is only an array in definition properties
          for (let i = 0; i < o['required'].length; i++) {
            // set required on the property if property key is included in required array
            set(o['properties'][o['required'][i]], 'required', true);
          }
        }

        // recursively apply the same (ie: convert desc to markdown)
        // operations on the object's nested structure
        this.deepRecursiveSwaggerUpdate(o);

      }

      // return the same object with modifications (eg: description to markdown)
      return o;

    });

  }

  /**
   * Get schema definition by a schema title (used by path.component.html)
   *
   * @param {string} title
   * @returns {Schema}
   * @memberof SwaggerService
   */
  getDefinition(title: string): Schema {

    const def = find(toArray(this._jsonSpec.definitions), (o: Schema) => {
      return o.title === title;
    });

    return def;

  }

  /**
   * Escape swagger definition path to valid 'id' html attribute
   * TODO: per chase, in the future we may want to have a global search
   * so that a participant can search content across our site from a single search
   * we'll need to consider how to architect this so that we can index various places
   * based on keywords ect. Somewhat related - https://github.com/GFTN/gftn-web/issues/73
   *
   * @param {string} ref - the id attribute on html element you want to jump to in the page
   * @memberof SwaggerService
   */
  jumpTo(kind: 'model' | 'path', operation?: 'get' | 'put' | 'post' | 'delete', title?: string) {

    // set the jumpText and navigate to location
    const jumpText = this.jumpText(kind, operation, title);

    // navigate
    this.router.navigate([
      '/docs/' +
      this.versionService.current.version +
      '/api/' + this.apiFileName
    ], { queryParams: { jump: jumpText } });

    return;

  }

  /**
   * generates text for and [id] attribute
   * text to identify which section in the page
   * to scroll (aka: jump) too
   * NOTE: used in conjunction with this.goTo()
   *
   * @param {('model' | 'path')} kind
   * @param {string} ref
   * @returns
   * @memberof SwaggerService
   */
  jumpText(kind: 'model' | 'path', operation?: string, title?: string) {

    // 'id' html attributes cannot have '#' or '/' or '{'
    // escape and transform special characters

    // sanitize the parameters if undefined
    let _kind = '';
    let _operation = '';
    let _title = '';
    if (kind) {
      _kind = kind;
    }
    if (operation) {
      _operation = operation;
    }
    if (title) {
      _title = title;
    }

    // create string with unescaped characters
    let jumpStr = _kind + '_' + _operation + '_' + _title;
    jumpStr = replace(jumpStr, /#/g, '_');
    jumpStr = replace(jumpStr, /\//g, '_');
    jumpStr = replace(jumpStr, /{/g, '_');
    jumpStr = replace(jumpStr, /}/g, '_');

    return jumpStr;

  }

  /**
   * Returns the last of item in #/definitions/xxxx
   *
   * @param {string} ref
   * @returns
   * @memberof SwaggerService
   */
  private refName(ref: string) {
    return last(split(ref, '/'));
  }

  /**
   * Convert text to markdown for inline docs (ie: descriptions in the docs)
   * IMPORTANT: Do not use in the angular template
   *
   * @param {string} mdText
   * @returns
   * @memberof SwaggerService
   */
  toMarkdown(mdText: string | number) {

    if (typeof (mdText) === 'number') {
      mdText = mdText.toString();
    }

    if (isString(mdText)) {

      // using a convention with '??' (ie:??SOMETEXT??) to replace text
      const _mdText = mdText.replace('??base_url??', location.origin)
        .replace('??version??', this.versionService.current.version);

      return this.markdownService.compile(_mdText);

      // return mdText;

    }

  }

}
