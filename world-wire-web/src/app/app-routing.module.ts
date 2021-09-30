// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { NgModule } from '@angular/core';
import { PageNotFoundComponent } from './not-found/not-found.component';
import { Routes } from '@angular/router/src/config';
import { RouterModule } from '@angular/router';
import { SecureVersionGuard } from './shared/guards/secure-version.guard';

const appRoutes: Routes = [
    {
        path: 'portal',
        loadChildren: './portal/portal.module#PortalModule'
    },
    {
        path: 'office',
        loadChildren: './office/office.module#OfficeModule'
    },
    {
        path: '',
        loadChildren: './public/public.module#PublicModule'
    },
    {
        path: 'soon',
        loadChildren: './coming/coming.module#ComingModule'
    },
    { path: 'not-found', component: PageNotFoundComponent },
    // redirect to most recent release
    {
        path: 'docs/secure',
        loadChildren: 'app/versions/secure-version/secure-version.module#VersionModule'
    },
    {
        path: 'docs', redirectTo: 'docs/v2.11.3', pathMatch: 'full',
        resolve: []
    },
    {
        path: 'docs/v2.11.3',

        loadChildren: 'app/versions/1/version.module#VersionModule',
        canActivate: [SecureVersionGuard],
    },
    {
        path: '2fa',
        loadChildren: './two-factor/two-factor.module#TwoFactorModule'
    },
    // ORDER MATTERS: catch all for not-found must come LAST!
    { path: '**', redirectTo: 'not-found' }

];

// appRoutes = VERSION_ROUTES().concat(appRoutes);

@NgModule({
    imports: [
        RouterModule.forRoot(
            // order of concat() matters, see above comment
            appRoutes,
            {
                // enableTracing: true, // <-- TODO: remove for production (debugging purposes only)
                scrollPositionRestoration: 'enabled'
            }
        )
    ],
    exports: [
        RouterModule
    ]
})
export class AppRoutingModule {

    // constructor(router: Router) {
    //     // unable to get lazy loaded routes to load and build dynamically
    //     router.resetConfig(concatRoutes);
    // }

}
