package apikey

import "time"

type Options struct {
	Pepper            []byte
	SecretBytes       int
	UsageQueueSize    int
	UsageFlushTimeout time.Duration
	MarkUsedTimeout   time.Duration
}

func (o Options) Normalize() Options {
	if o.SecretBytes == 0 {
		o.SecretBytes = DefaultSecretBytes
	}
	if o.UsageQueueSize == 0 {
		o.UsageQueueSize = 1024
	}
	if o.UsageFlushTimeout == 0 {
		o.UsageFlushTimeout = 250 * time.Millisecond
	}
	if o.MarkUsedTimeout == 0 {
		o.MarkUsedTimeout = 500 * time.Millisecond
	}
	return o
}

func (o Options) Validate() error {
	o = o.Normalize()
	if len(o.Pepper) < MinSecretBytes {
		return ErrInvalidPepper
	}
	if o.SecretBytes < MinSecretBytes {
		return ErrInvalidSecretBytes
	}
	if o.UsageQueueSize < 1 {
		return ErrInvalidUsageQueueSize
	}
	if o.UsageFlushTimeout <= 0 || o.MarkUsedTimeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}
