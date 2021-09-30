// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { Request, Response, NextFunction } from 'express';
// mock middleware for testing purpose; not implemented in else where for now
export class MockMiddleWare {
    authenticateIbmId = async (req: Request, res: Response, next: NextFunction) => {
        req.user = {_json:{email : "your.user@your.domain"}};
        req.user.emailAddress = "your.user@your.domain";
        console.debug(req.user);
        next();
    }
}