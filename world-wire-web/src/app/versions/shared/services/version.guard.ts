// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { SwaggerService } from './swagger.service';
import { get, findLast } from 'lodash';
import { VERSION_DETAILS } from '../../../shared/constants/versions.constant';
import { VersionService } from './version.service';


@Injectable()
export class VersionGuard implements CanActivate {
    constructor(
        private swaggerService: SwaggerService,
        private versionService: VersionService,
        private router: Router
    ) { }

    canActivate(
        route: ActivatedRouteSnapshot,
        state: RouterStateSnapshot
    ): Observable<boolean> | Promise<boolean> | boolean {

        // fallback to window.location if router doesn't work
        const currentUrl = state.url ? state.url : window.location.pathname;

        if (currentUrl.split('/')[1] === 'docs') {
            const routeVersion = currentUrl.split('/')[2];

            // onnly set current if not set
            if (!this.versionService.current || routeVersion !== this.versionService.current.version) {
                // set the 'current' version by the version defined in the route
                this.versionService.current = this.versionService.getVersion(routeVersion);

                // could not find version details. version does not exist
                if (!this.versionService.current) {
                    this.router.navigate(['/not-found']);
                }
            }
        } else {
            // set version to newest release if the url does not define the route
            this.versionService.current = this.versionService.current ? this.versionService.current : this.versionService.getNewestVersion();
        }

        // reset swagger service
        this.swaggerService._jsonSpec = null;
        this.swaggerService.apiFileName = null;
        this.swaggerService.navPaths = null;

        // get api fileName for route
        const apiFileName = get(route.params, 'name');

        // if no apiFileName
        if (!apiFileName) {
            return true;
        } else {
            // get swagger file over https
            return new Promise((resolve, reject) => {

                // https get swagger json definition
                this.swaggerService.populateSwaggerDocs(apiFileName)
                    .then(() => {
                        // console.log('done last');
                        resolve(true);
                    });

            });
        }

    }
}
