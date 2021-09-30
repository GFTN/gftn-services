// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { trigger, style, animate, transition, state } from '@angular/animations';

// Animation Documentation
// http://www.carbondesignsystem.com/guidelines/motion/#easing
// https://angular.io/guide/animations

// the fade-in/fade-out animation
export const simpleFadeAnimation = trigger('simpleFadeAnimation', [

    // the "in" style determines the "resting" state of the element when it is visible.
    state('in', style({
        opacity: 1
    })),

    // fade in when created. this could also be written as transition('void => *')
    transition(':enter', [
        style({
            opacity: 0
        }),
        animate('300ms cubic-bezier(0.5, 0, 0.1, 1)')
    ]),

    // fade out when destroyed. this could also be written as transition('void => *')
    transition(':leave',
        animate('300ms cubic-bezier(0.5, 0, 0.1, 1)',
            style({
                opacity: 0
            }))
    )

]);

// the slide-in/slide-out animation
export const slideAnimation = trigger(
    'slideAnimation', [
        // slide out and fade in
        transition(':enter', [
            style({
                transform: 'translateX(100%)',
                opacity: 0
            }),
            // Carbon Design's Motion easing
            animate('300ms cubic-bezier(0.5, 0, 0.1, 1)',
                style({
                    transform: 'translateX(0)',
                    opacity: 1
                }))
        ]),
        // slide in and fade out
        transition(':leave', [
            style({
                transform: 'translateX(0)',
                opacity: 1
            }),
            // Carbon Design's Motion easing
            animate('300ms cubic-bezier(0.5, 0, 0.1, 1)',
                style({
                    transform: 'translateX(100%)',
                    opacity: 0
                }))
        ])
    ]
);

// slide down/up animation. Used for menus, popovers, etc.
export const slideDownAnimation = trigger(
    'slideDownAnimation', [
        // the "in" style determines the "resting" state of the element when it is visible.
        state('in',
            style({
                height: 'auto',
                opacity: 1,
                overflow: 'visible'
            })
        ),

        // slide in when created. this could also be written as transition('void => *')
        transition(':enter', [
            style({
                height: 0,
                opacity: 0,
                overflow: 'hidden'
            }),
            animate('300ms cubic-bezier(0.5, 0, 0.1, 1)',
                style({
                    height: 'auto',
                    opacity: 1,
                    overflow: 'visible'
                })
            )
        ]),

        // slide out when destroyed. this could also be written as transition('void => *')
        transition(':leave', [
            animate('300ms cubic-bezier(0.5, 0, 0.1, 1)',
                style({
                    height: 0,
                    opacity: 0,
                    overflow: 'hidden'
                })
            )
        ])

    ]
);
