-- SQL function for updateQuote
-- RFI allows to update quote if : statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_CANCELLED OR statusQuote = _QUOTESTATUS_FAILED
CREATE OR REPLACE FUNCTION UpdateQuote(
		_quoteID varchar,
 		_rfiID varchar, 
		_statusQuote int, 
		_exchangeRate numeric , 
		_timequote int, 
		_timeStartRfi int , 
		_timeExpireRfi int, 
		_addressReceiveRfi varchar, 
		_addressSendRfi varchar, 
		_quoteResponse jsonb, 
		_quoteResponseBase64 varchar, 
		_quoteResponseSignature varchar, 
		_limitMaxRfi numeric, 
		_limitMinRfi numeric, 
		_sourceAsset jsonb, 
		_targetAsset jsonb, 
		_limitMaxOfi numeric, 
		_limitMinOfi numeric, 
		_timeExpireOFI bigint, 
		_ofiID varchar,
		_QUOTESTATUS_EXECUTED int,
		_QUOTESTATUS_CANCELLED int,
		_QUOTESTATUS_EXECUTING int,
		_QUOTESTATUS_PENDING int,
		_QUOTESTATUS_UPDATED int,
		_QUOTESTATUS_FAILED int
		) RETURNS INT
	AS $$
	DECLARE
	statusQuote int;
	sourceAsset jsonB;
	targetAsset jsonB;
	timeExpireOFI bigint;
	limitMaxOfi numeric;
	limitMinOfi numeric;
	ofiID varchar;
	rowCount int;
	BEGIN
	SELECT q.statusQuote, q.sourceAsset, q.targetAsset, q.timeExpireOFI, q.limitMaxOfi, q.limitMinOfi, q.ofiID into statusQuote, sourceAsset, targetAsset, timeExpireOFI, limitMaxOfi, limitMinOfi, ofiID
	FROM quotes q WHERE q.quoteID = _quoteID AND q.rfiID = _rfiID FOR UPDATE;
	GET DIAGNOSTICS rowCount = ROW_COUNT;
	IF rowCount = 0 THEN
	RAISE EXCEPTION 'No record found';
	END IF;
	IF sourceAsset != _sourceAsset THEN
	RAISE EXCEPTION 'Source asset mismatches quote request';
	END IF;
	IF targetAsset != _targetAsset THEN
	RAISE EXCEPTION 'Target asset mismatches quote request';
	END IF;
	IF ofiID != _ofiID THEN
	RAISE EXCEPTION 'OFIID mismatches quote request';
	END IF;
	IF limitMaxOfi != _limitMaxOfi OR limitMinOfi != _limitMinOfi THEN
	RAISE EXCEPTION 'limit range mismatches quote request';
	END IF;
	IF timeExpireOFI != _timeExpireOFI THEN
	RAISE EXCEPTION 'timeExpireOFI Mismatches quote request';
	END IF;
	IF rowCount > 1 THEN
	RAISE EXCEPTION 'Multiple record found';
	END IF;
	
	IF statusQuote = _QUOTESTATUS_EXECUTING THEN
	RAISE EXCEPTION 'Quote is executing';
	END IF;
	IF statusQuote = _QUOTESTATUS_EXECUTED THEN
	RAISE EXCEPTION 'Quote is executed';
	END IF;
	IF statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_CANCELLED OR statusQuote = _QUOTESTATUS_FAILED THEN
	UPDATE quotes q
	SET statusQuote = _statusQuote, exchangeRate = _exchangeRate, timequote = _timequote, timeStartRfi = _timeStartRfi, timeExpireRfi = _timeExpireRfi,
	AddressReceiveRfi = _addressReceiveRfi, AddressSendRfi = _addressSendRfi, QuoteResponse = _quoteResponse, QuoteResponseBase64 = _quoteResponseBase64, QuoteResponseSignature = _quoteResponseSignature, limitMaxRfi = _limitMaxRfi, limitMinRfi = _limitMinRfi
	WHERE q.quoteID = _quoteID AND q.rfiID = _rfiID;
	END IF;
	RETURN statusquote;

	END;
	$$ LANGUAGE plpgsql;


-- SQL function for executingQuote
-- OFI allow to execute quote if: statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_FAILED
CREATE OR REPLACE FUNCTION ExecutingQuote(
		_quoteID varchar, 
		_ofiID varchar, 
		_QuoteResponse jsonb, 
		_QUOTESTATUS_EXECUTED int,
		_QUOTESTATUS_CANCELLED int,
		_QUOTESTATUS_EXECUTING int,
		_QUOTESTATUS_PENDING int,
		_QUOTESTATUS_UPDATED int,
		_QUOTESTATUS_FAILED int,
		_timeExecuting bigint, 
		_amount numeric
		) RETURNS INT
	AS $$
	DECLARE
	statusQuote int;
	quoteResponse jsonB;
	timeExpireRFI bigint;
	limitMaxRfi numeric;
	limitMinRfi numeric;
	rowCount int;
	BEGIN
	SELECT q.statusQuote, q.quoteResponse, q.timeExpireRFI, q.limitMaxRfi, q.limitMinRfi into statusQuote, quoteResponse, timeExpireRFI, limitMaxRfi, limitMinRfi
	FROM quotes q WHERE q.quoteID = _quoteID AND q.ofiID = _ofiID FOR UPDATE;
	GET DIAGNOSTICS rowCount = ROW_COUNT;
	IF rowCount = 0 THEN
	RAISE EXCEPTION 'No record found';
	END IF;
	IF quoteResponse != _QuoteResponse THEN
	RAISE EXCEPTION 'quote Response mismatches ';
	END IF;
	IF limitMaxRfi < _amount OR limitMinRfi > _amount THEN
	RAISE EXCEPTION 'amount falls out of RFI''s range limits';
	END IF;
	IF timeExpireRFI < _timeExecuting THEN
	RAISE EXCEPTION 'quote expired';
	END IF;
	IF rowCount > 1 THEN
	RAISE EXCEPTION 'Multiple record found';
	END IF;
	
	IF statusQuote = _QUOTESTATUS_EXECUTING THEN
	RAISE EXCEPTION 'Quote is executing';
	END IF;
	IF statusQuote = _QUOTESTATUS_EXECUTED THEN
	RAISE EXCEPTION 'Quote is executed';
	END IF;
	IF statusQuote = _QUOTESTATUS_CANCELLED THEN
	RAISE EXCEPTION 'Quote is cancelled';
	END IF;
	IF statusQuote = _QUOTESTATUS_PENDING THEN
	RAISE EXCEPTION 'Quote is pending for update';
	END IF;
	IF statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_FAILED THEN
	UPDATE quotes q
	SET statusQuote = _QUOTESTATUS_EXECUTING, amount = _amount
	WHERE q.quoteID = _quoteID AND q.ofiID = _ofiID AND q.quoteResponse = _QuoteResponse;
	END IF;
	RETURN statusquote;

	END;
	$$ LANGUAGE plpgsql;

-- SQL function for cancelQuote
-- RFI allows to cancel quote if: statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_CANCELLED OR statusQuote = _QUOTESTATUS_FAILED THEN
CREATE OR REPLACE FUNCTION CancelQuote(
		_quoteID varchar,
		_rfiID varchar,
		_QUOTESTATUS_CANCELLED int,
		_timeCancel bigint,
		_QUOTESTATUS_UPDATED int,
		_QUOTESTATUS_FAILED int,
		_QUOTESTATUS_PENDING int,
		_QUOTESTATUS_EXECUTED int,
		_QUOTESTATUS_EXECUTING int
		) RETURNS INT
	AS $$
	DECLARE
	statusQuote int;
	rowCount int;
	BEGIN
	SELECT q.statusQuote into statusQuote
	FROM quotes q WHERE q.quoteID = _quoteID AND q.rfiID = _rfiID FOR UPDATE;
	GET DIAGNOSTICS rowCount = ROW_COUNT;
	IF rowCount = 0 THEN
	RAISE EXCEPTION 'No record found';
	END IF;
	IF rowCount > 1 THEN
	RAISE EXCEPTION 'Multiple record found';
	END IF;
	IF statusQuote = _QUOTESTATUS_EXECUTING THEN
	RAISE EXCEPTION 'Quote is executing';
	END IF;
	IF statusQuote = _QUOTESTATUS_EXECUTED THEN
	RAISE EXCEPTION 'Quote is executed';
	END IF;
	IF statusQuote = _QUOTESTATUS_PENDING THEN
	RAISE EXCEPTION 'Quote is pending for update';
	END IF;
	IF statusQuote <= _QUOTESTATUS_UPDATED OR statusQuote = _QUOTESTATUS_CANCELLED OR statusQuote = _QUOTESTATUS_FAILED THEN
	UPDATE quotes
	SET statusQuote = _QUOTESTATUS_CANCELLED, timeCancel = _timeCancel
	WHERE quoteID = _quoteID AND rfiID = _rfiID;
	END IF;
	RETURN statusquote;

	END;
	$$ LANGUAGE plpgsql;