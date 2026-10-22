package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestRetryRepositoryRead(t *testing.T) {
	t.Run("returns immediately on success", func(t *testing.T) {
		calls := 0
		err := retryRepositoryRead(func(*schema.ResourceData, interface{}) error {
			calls++
			return nil
		}, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 1, calls)
	})

	t.Run("retries until success", func(t *testing.T) {
		oldDelay := repositoryReadRetryDelay
		repositoryReadRetryDelay = time.Millisecond
		t.Cleanup(func() { repositoryReadRetryDelay = oldDelay })

		calls := 0
		err := retryRepositoryRead(func(*schema.ResourceData, interface{}) error {
			calls++
			if calls < 3 {
				return errors.New("not ready")
			}
			return nil
		}, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, 3, calls)
	})

	t.Run("retries while a written repository is not yet found", func(t *testing.T) {
		oldDelay := repositoryReadRetryDelay
		repositoryReadRetryDelay = time.Millisecond
		t.Cleanup(func() { repositoryReadRetryDelay = oldDelay })

		resourceData := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
		resourceData.SetId("example")
		calls := 0
		err := retryRepositoryRead(func(d *schema.ResourceData, _ interface{}) error {
			calls++
			if calls < 3 {
				d.SetId("")
			}
			return nil
		}, resourceData, nil)
		assert.NoError(t, err)
		assert.Equal(t, 3, calls)
		assert.Equal(t, "example", resourceData.Id())
	})

	t.Run("returns last error after max attempts", func(t *testing.T) {
		oldDelay := repositoryReadRetryDelay
		repositoryReadRetryDelay = time.Millisecond
		t.Cleanup(func() { repositoryReadRetryDelay = oldDelay })

		expectedErr := errors.New("persistent failure")
		calls := 0
		err := retryRepositoryRead(func(*schema.ResourceData, interface{}) error {
			calls++
			return expectedErr
		}, nil, nil)
		assert.ErrorIs(t, err, expectedErr)
		assert.Equal(t, repositoryReadMaxAttempts, calls)
	})
}
