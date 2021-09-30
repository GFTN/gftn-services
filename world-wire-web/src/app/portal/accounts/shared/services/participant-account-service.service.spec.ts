// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { TestBed, inject } from '@angular/core/testing';
import {
  HttpClientTestingModule,
} from '@angular/common/http/testing';
import { RouterTestingModule } from '@angular/router/testing';

import { ParticipantAccountService as ParticipantAccountService } from './participant-account-service.service';
import { ParticipantAccount } from '../../../../shared/models/account.interface';
import { SessionService } from '../../../../shared/services/session.service';
import { AuthService } from '../../../../shared/services/auth.service';
import { AngularFireAuth } from '@angular/fire/auth';
import { of } from 'rxjs';
import { DocumentService } from '../../../../shared/services/document.service';

import { AngularFireModule } from '@angular/fire';
import { AngularFireDatabaseModule, AngularFireDatabase } from '@angular/fire/database';
import { environment } from '../../../../../environments/environment.dev';
import { HttpHeaders } from '@angular/common/http';
import { KillSwitchRequestDetail } from '../../../shared/models/killswitch-request.interface';

describe('ParticipantAccountService', () => {
  let service: ParticipantAccountService;

  const account: ParticipantAccount = {
    address: 'testaddress12345',
  };

  let authService: AuthService;

  let sessionService: SessionService;

  let db: AngularFireDatabase;

  const authState: any = {
    displayName: null,
    isAnonymous: true,
    email: 'test@email.com',
    // uid: '17WvU2Vj58SnTz8v7EqyYYb0WRc2'
    uid: '1234567890',
    // delete: () => Promise.resolve(),
    getIdToken: () => Promise.resolve('1234567890'),
    // getIdTokenResult: () => Promise.resolve(null),
    // emailVerified: true,
  };

  const mockAngularFireAuth: any = {
    auth: jasmine.createSpyObj('auth', {
      'signInAnonymously': Promise.resolve({
        code: 'auth/operation-not-allowed'
      }),
      'onAuthStateChanged': (callback) => {
        callback({
          user: authState
        });
      },
      'currentUser': authState,
      // 'signInWithPopup': Promise.reject(),
      // 'signOut': Promise.reject()
    }),
    authState: of(authState)
  };

  const mockSessionService: any = {
    institution: {
      info: {
        institutionId: 'testInstitutionId1',
        slug: 'big-test-bank',
        status: 'active',
        name: 'Big Test Bank',
        country: 'USA',
        address1: '600 Anton',
        address2: '',
        city: 'Costa Mesa',
        state: 'CA',
        zip: '92626',
        kind: 'Bank'
      },
    },
    currentNode: {
      institutionId: 'testInstitutionId1',
      status: ['complete'],
      bic: 'XXXXXX',
      role: 'MM',
      participantId: 'unittest1',
      initialized: true,
      countryCode: 'USA',
    }
  };

  const mockAuthService: any = {
    auth: mockAngularFireAuth,
    getFirebaseIdToken: () => {

      const h = new HttpHeaders();

      h.set('x-fid', 'testfid12345');
      h.set('x-iid', 'testInstitutionId1');

      Promise.resolve(h);
    },
    addMakerCheckerHeaders: () => {

      const h = new HttpHeaders();

      h.set('x-permission', 'request');

      return h;
    },
  };

  let mockSuspendReactivateRequest: KillSwitchRequestDetail;

  const user1 = 'user1@test.com';

  beforeEach(() => {

    TestBed.configureTestingModule({
      providers: [
        { provide: AngularFireAuth, useValue: mockAngularFireAuth },
        // { provide: AngularFireDatabase, useValue: afDbMock },
        DocumentService,
        { provide: AuthService, useValue: mockAuthService },
        { provide: SessionService, useValue: mockSessionService },
        ParticipantAccountService,
      ],
      imports: [
        AngularFireModule.initializeApp(environment.firebase),
        AngularFireDatabaseModule,
        HttpClientTestingModule,
        RouterTestingModule
      ],
    });
    inject([AngularFireDatabase], (_db: AngularFireDatabase) => {
      db = _db;

      db.database.goOffline();
    })();

    sessionService = TestBed.get(SessionService);

    authService = TestBed.get(AuthService);

    service = TestBed.get(ParticipantAccountService);

    // set up mock request for testing approval/rejection/getting data
    mockSuspendReactivateRequest = {
      key: account.address,
      participantId: sessionService.currentNode.participantId,
      accountAddress: account.address,
      approvalIds: ['unittestApprovalId1'],
      status: 'suspend_requested',
    };
  });

  it('should be created', () => {
    service = TestBed.get(ParticipantAccountService);
    expect(service).toBeTruthy();
  });

  it('#getSupendReactivateRequest() should set up db ref properly', () => {
    const result = service.getSupendReactivateRequest(account);

    // request ref should not be null
    expect(service.suspendReactivateDbRef).toBeTruthy();
  });


  it('#requestSuspendReactivateAccount() should fail gracefully without proper setup', () => {

    setTimeout(async () => {

      // fallback should return null if db is not set up
      const result1 = await service.requestSuspendReactivateAccount(account, false);

      expect(result1).toBeNull();
    });

  });

  it('should run #requestSuspendReactivateAccount() to request suspension of account', () => {

    setTimeout(async () => {

      // setting up actual request
      await service.getSupendReactivateRequest(account);

      const result2 = await service.requestSuspendReactivateAccount(account, true);

      expect(result2).toBeTruthy();
    });

  });

  it('should run #requestSuspendReactivateAccount() to request reactivation of account', () => {

    setTimeout(async () => {
      // setting up actual request
      await service.getSupendReactivateRequest(account);

      const result2 = await service.requestSuspendReactivateAccount(account, false);

      // request should return an actual result
      expect(result2).toBeTruthy();
    });

  });

  it('#approveSuspendReactivateAccount() should fail gracefully without request or proper setup', () => {

    setTimeout(async () => {

      const result1 = await service.approveSuspendReactivateAccount(account);

      // fallback should return null if db is not set up
      expect(result1).toBeNull();

      await service.getSupendReactivateRequest(account);

      const result2 = await service.approveSuspendReactivateAccount(account);

      // fallback should return null if no request to approve
      expect(result2).toBeNull();
    });
  });

  it('should run #approveSuspendReactivateAccount() to approve account suspension request', () => {

    setTimeout(async () => {

      // setting up actual request
      await service.getSupendReactivateRequest(account);

      service.suspendReactivateRequest = mockSuspendReactivateRequest;

      // approval should return an actual result
      const result2 = await service.approveSuspendReactivateAccount(account);

      expect(result2).toBeTruthy();
    });
  });

  it('#rejectSuspendReactivateAccount() should fail gracefully without request or proper setup', () => {

    setTimeout(async () => {

      // fallback should return null
      const result1 = await service.rejectSuspendReactivateAccount(account);

      expect(result1).toBeNull();

      await service.getSupendReactivateRequest(account);

      // fallback should return null if no request to reject
      const result2 = await service.rejectSuspendReactivateAccount(account);

      expect(result2).toBeNull();
    });
  });

  it('should run #rejectSuspendReactivateAccount() to reject account suspension request', () => {

    setTimeout(async () => {

      // setting up actual request
      await service.getSupendReactivateRequest(account);

      service.suspendReactivateRequest = mockSuspendReactivateRequest;

      // rejection should return an actual result
      const result2 = await service.rejectSuspendReactivateAccount(account);

      expect(result2).toBeTruthy();
    });
  });

  it('#getAccountStatus() should always return a value', async () => {
    const result = await service.getAccountStatus(account.address);

    // default status for no killswitch request
    expect(result).toBeTruthy();

    expect(result).toEqual('normal');
  });

  it('#getAccountStatus() should return the proper status of existing request', async () => {

    service.suspendReactivateRequest = mockSuspendReactivateRequest;
    service.suspendReactivateRequest.status = 'suspend_requested';

    let result = await service.getAccountStatus(account.address);

    // status should be 'suspend_requested'
    expect(result).toEqual('suspend_requested');

    service.suspendReactivateRequest.status = 'reactivate_requested';

    result = await service.getAccountStatus(account.address);

    // status should be 'reactivate_requested'
    expect(result).toEqual('reactivate_requested');

    service.suspendReactivateRequest.status = 'suspended';

    result = await service.getAccountStatus(account.address);

    // status should be 'suspended'
    expect(result).toEqual('suspended');

    service.suspendReactivateRequest.status = 'normal';

    result = await service.getAccountStatus(account.address);

    // status should be 'normal'
    expect(result).toEqual('normal');

  });

  it('should get email of requestor of suspend request', async () => {

    // Suspension Requested
    service.suspendReactivateRequest = mockSuspendReactivateRequest;
    service.suspendReactivateRequest.suspendRequestedBy = user1;
    service.suspendReactivateRequest.status = 'suspend_requested';

    const result = await service.getRequesterEmail(account);

    expect(result).toEqual(user1);
  });

  it('should get email of requestor of reactivate request', async () => {

    // Reactivation Requested
    service.suspendReactivateRequest = mockSuspendReactivateRequest;
    service.suspendReactivateRequest.reactivateRequestedBy = user1;
    service.suspendReactivateRequest.status = 'reactivate_requested';

    const result = await service.getRequesterEmail(account);

    expect(result).toEqual(user1);
  });

  it('should get suspension action of request', async () => {

    service.suspendReactivateRequest = mockSuspendReactivateRequest;
    service.suspendReactivateRequest.status = 'suspend_requested';

    const result = await service.getSuspendReactivate(account);

    expect(result).toEqual('Suspension');
  });

  it('should get reactivation action of request', async () => {

    service.suspendReactivateRequest = mockSuspendReactivateRequest;
    service.suspendReactivateRequest.status = 'reactivate_requested';

    const result = await service.getSuspendReactivate(account);

    expect(result).toEqual('Reactivation');
  });

  it('should return no action if no request', async () => {

    const result = await service.getSuspendReactivate(account);

    expect(result).toEqual(null);
  });
});
