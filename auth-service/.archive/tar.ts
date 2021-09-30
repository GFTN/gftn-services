// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as tar from 'tar';
// import * as tar2 from 'tar-stream';
// import * as fs from 'fs'
// import * as zlib from 'zlib' 

export class Tar {

    /**
     * tar files and compress 
     *
     * @param {string} inputDirPath
     * @param {string} outputFilename
     * @memberof Compress
     */
    async compress(inputDirPath: string, outputFilename: string){

        await tar.c(
            {
              gzip: true,
              file: outputFilename
            },
            [inputDirPath]
          ) 

    }

    /**
     * decompress and extract files 
     *
     * @param {string} inputDirPath
     * @param {string} outputFilename
     * @memberof Compress
     */
    async extract(inputFilename: string){

        await tar.x(
            {
                file: inputFilename
            }
        )

    }

    // async extractStream() {

    //     return new Promise((resolve, reject) => {

    //         var extract = tar2.extract()

    //         extract.on('entry', function (header, stream, callback) {
                
    //             // header is the tar header
    //             // stream is the content body (might be an empty stream)
    //             // call next when you are done with this entry
    //             console.log(header);

    //             let data='';

    //             stream.on('data', (chunk)=> {
    //                 data += chunk;
    //                 console.log(data);
    //             });
                
    //             stream.on('end', function () {
    //                 console.log('end');
    //                 callback() // ready for next entry
    //             })

    //             stream.resume() // just auto drain the stream

    //         })

    //         extract.on('finish', function () {
    //             console.log('finish');
    //             resolve();
    //         })

    //         // const plainTar = fs.readFileSync('testtar.tar')
    //         // fs.createReadStream(plainTar).pipe(process.stdout);

    //         // extract.end(fs.readFileSync('testtar.tar'))
    //         // //   extract.end(fs.readFileSync('testtar.tgz'))

    //         // const gzippedStream = fs.createReadStream('testtar.tgz').pipe(gunzip())
    //         // gzippedStream.pipe(process.stdout);

    //         // extract.end(gzippedStream)

    //         fs.createReadStream('testtar.tgz')
    //         .pipe(zlib.createGunzip())
    //         .pipe(extract);

    //     });

    // }

}