-- Quotes schema; denormalized for active query
CREATE TABLE QUOTES(
	QuoteID VARCHAR(100) NOT NULL UNIQUE ,
	RequestID VARCHAR(100) NOT NULL, 
	RfiID VARCHAR(100) NOT NULL, 
	OfiID VARCHAR(100) NOT NULL,
	LimitMaxOfi NUMERIC(25,10) NOT NULL CHECK (LimitMinOfi <= LimitMaxOfi),
	LimitMinOfi NUMERIC(25,10) NOT NULL CHECK (LimitMinOfi <= LimitMaxOfi),
	LimitMaxRfi NUMERIC(25,10) CHECK (LimitMinOfi <= LimitMaxOfi),
	LimitMinRfi NUMERIC(25,10) CHECK (LimitMinRfi <= LimitMaxRfi),
	Amount NUMERIC(25,10) CHECK (Amount <= LimitMaxRfi) CHECK (Amount >= LimitMinRfi), -- source amount upon execution
	ExchangeRate NUMERIC(25,10),
	SourceAsset JSONB NOT NULL, 
	TargetAsset JSONB NOT NULL, 
	TimeRequest BIGINT NOT NULL, 
	TimeQuote BIGINT, 
	StatusQuote INT NOT NULL, 
	TimeExpireOfi BIGINT,
	TimeStartRfi BIGINT CHECK (TimeStartRfi <= TimeExpireRfi),
	TimeExpireRfi BIGINT CHECK (TimeStartRfi <= TimeExpireRfi),
	TimeExecuted BIGINT, 
	TimeCancel BIGINT,
	AddressReceiveRfi VARCHAR(100),
	AddressSendRfi VARCHAR(100),
	QuoteResponse JSONB,
	QuoteResponseBase64 VARCHAR(2000),
	QuoteResponseSignature VARCHAR(1000)
	);    

-- Request schema
CREATE TABLE REQUESTS (
	RequestID  VARCHAR(100) NOT NULL UNIQUE,
	TimeExpireOfi BIGINT,
	LimitMaxOfi NUMERIC(25,10) NOT NULL CHECK (LimitMinOfi <= LimitMaxOfi),
	LimitMinOfi NUMERIC(25,10) NOT NULL CHECK (LimitMinOfi <= LimitMaxOfi),
	SourceAsset JSONB NOT NULL,  
	TargetAsset JSONB NOT NULL, 
	TimeRequest BIGINT NOT NULL, 
	OfiID VARCHAR(100) NOT NULL
	);
	
-- index for quotes 
-- postgre create index for unique field automatically
CREATE INDEX idx_quotes_rfiID_statusquote
ON quotes(rfiID, statusquote);
CREATE INDEX idx_quotes_ofiID_statusquote
ON quotes(ofiID, statusquote);