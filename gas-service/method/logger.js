// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var colors = require('colors');

module.exports = function (groupname){
    this.groupname = groupname
    this.logger = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string))
    }
    this.info = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).green)
    }
    this.error = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).red)
    }
    this.silly = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).rainbow)
    }
    this.input = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).gray)
    }
    this.verbose = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).grey)
    }
    this.prompt = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).grey)
    }
    this.data = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).grey)
    }
    this.warn = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).yellow)
    }
    this.debug = function  (tittle,string){
        console.log((new Date().toISOString() + ' '+this.groupname+ '\t \u25ba ' + tittle + ' : ' + string).blue)
    }
}
