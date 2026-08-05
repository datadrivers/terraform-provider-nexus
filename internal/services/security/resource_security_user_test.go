package security_test

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/datadrivers/go-nexus-client/nexus3/schema/security"
	"github.com/datadrivers/terraform-provider-nexus/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccResourceSecurityUser() security.User {
	return security.User{
		UserID:       fmt.Sprintf("user-test-%s", acctest.RandString(10)),
		FirstName:    fmt.Sprintf("user-firstname-%s", acctest.RandString(10)),
		LastName:     fmt.Sprintf("user-lastname-%s", acctest.RandString(10)),
		EmailAddress: fmt.Sprintf("user-email-%s@example.com", acctest.RandString(10)),
		Status:       "active",
		Password:     acctest.RandString(16),
		Roles:        []string{"nx-admin"},
		Source:       "default",
	}
}

func TestAccResourceSecurityUser(t *testing.T) {
	resName := "nexus_security_user.acceptance"

	user := testAccResourceSecurityUser()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSecurityUserConfig(user),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr(resName, "id", user.UserID),
					resource.TestCheckResourceAttr(resName, "userid", user.UserID),
					resource.TestCheckResourceAttr(resName, "firstname", user.FirstName),
					resource.TestCheckResourceAttr(resName, "lastname", user.LastName),
					resource.TestCheckResourceAttr(resName, "password", user.Password),
					resource.TestCheckResourceAttr(resName, "email", user.EmailAddress),
					resource.TestCheckResourceAttr(resName, "status", user.Status),
					resource.TestCheckResourceAttr(resName, "roles.#", strconv.Itoa(len(user.Roles))),
					resource.TestCheckResourceAttr(resName, "source", user.Source),
					// FIXME: (BUG) Incorrect roles state representation.
					// For some reasons, 1st element in array is not stored as roles.0, but instead it's stored
					// as roles.3360874991 where 3360874991 is a "random" number.
					// This number changes from test run to test run.
					// It may be a pointer to int instead of int itself, but it's not clear and requires additional research.
					// resource.TestCheckResourceAttr(resName, "roles.3360874991", "nx-admin"),
				),
			},
			{
				ResourceName:      resName,
				ImportStateId:     user.UserID,
				ImportState:       true,
				ImportStateVerify: true,
				// Password is not returned
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

// Test write-only password functionality
func TestAccResourceSecurityUser_WriteOnlyPassword(t *testing.T) {
	resName := "nexus_security_user.acceptance"

	user := testAccResourceSecurityUser()
	updatedFirstName := fmt.Sprintf("user-firstname-updated-%s", acctest.RandString(10))
	newPassword := acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			// Create user with write-only password
			{
				Config: testAccResourceSecurityUserWriteOnlyConfig(user, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "id", user.UserID),
					resource.TestCheckResourceAttr(resName, "userid", user.UserID),
					resource.TestCheckResourceAttr(resName, "firstname", user.FirstName),
					resource.TestCheckResourceAttr(resName, "lastname", user.LastName),
					resource.TestCheckResourceAttr(resName, "email", user.EmailAddress),
					resource.TestCheckResourceAttr(resName, "status", user.Status),
					resource.TestCheckResourceAttr(resName, "password_wo_version", "1"),
					resource.TestCheckResourceAttr(resName, "roles.#", strconv.Itoa(len(user.Roles))),
					// password_wo should NOT be in state
					resource.TestCheckNoResourceAttr(resName, "password_wo"),
					// legacy password should NOT be set
					resource.TestCheckNoResourceAttr(resName, "password"),
				),
			},
			// Update non-password fields (should not trigger password change)
			{
				Config: testAccResourceSecurityUserWriteOnlyConfigWithUpdatedName(user, updatedFirstName, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "firstname", updatedFirstName),
					resource.TestCheckResourceAttr(resName, "password_wo_version", "1"), // Version unchanged
				),
			},
			// Update password by changing version
			{
				Config: testAccResourceSecurityUserWriteOnlyConfigWithPasswordUpdate(user, updatedFirstName, newPassword, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resName, "firstname", updatedFirstName),
					resource.TestCheckResourceAttr(resName, "password_wo_version", "2"), // Version incremented
					// password_wo should still NOT be in state
					resource.TestCheckNoResourceAttr(resName, "password_wo"),
				),
			},
			// Import test for write-only password
			{
				ResourceName:      resName,
				ImportStateId:     user.UserID,
				ImportState:       true,
				ImportStateVerify: true,
				// password_wo is write-only and should be ignored
				ImportStateVerifyIgnore: []string{"password_wo"},
			},
		},
	})
}

// Test that legacy and write-only password fields conflict
func TestAccResourceSecurityUser_PasswordConflict(t *testing.T) {
	user := testAccResourceSecurityUser()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSecurityUserConflictConfig(user),
				ExpectError: regexp.MustCompile("conflicts with password"),
			},
		},
	})
}

// Test that password_wo requires password_wo_version
func TestAccResourceSecurityUser_WriteOnlyPasswordRequiresVersion(t *testing.T) {
	user := testAccResourceSecurityUser()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acceptance.AccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceSecurityUserWriteOnlyNoVersionConfig(user),
				ExpectError: regexp.MustCompile("all of `password_wo,password_wo_version` must be specified"), // ← Исправленный regex
			},
		},
	})
}

// Legacy password configuration (backward compatibility)
func testAccResourceSecurityUserConfig(user security.User) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid    = "%s"
	firstname = "%s"
	lastname  = "%s"
	email     = "%s"
	password  = "%s"
	status    = "%s"
	roles     = ["%s"]
	source    = "%s"
}
`, user.UserID, user.FirstName, user.LastName, user.EmailAddress, user.Password, user.Status, strings.Join(user.Roles, "\", \""), user.Source)
}

// Write-only password configuration
func testAccResourceSecurityUserWriteOnlyConfig(user security.User, version int) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid              = "%s"
	firstname           = "%s"
	lastname            = "%s"
	email               = "%s"
	password_wo         = "%s"
	password_wo_version = %d
	status              = "%s"
	roles               = ["%s"]
}
`, user.UserID, user.FirstName, user.LastName, user.EmailAddress, user.Password, version, user.Status, strings.Join(user.Roles, "\", \""))
}

// Write-only password configuration with updated firstname
func testAccResourceSecurityUserWriteOnlyConfigWithUpdatedName(user security.User, updatedFirstName string, version int) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid              = "%s"
	firstname           = "%s"
	lastname            = "%s"
	email               = "%s"
	password_wo         = "%s"
	password_wo_version = %d
	status              = "%s"
	roles               = ["%s"]
}
`, user.UserID, updatedFirstName, user.LastName, user.EmailAddress, user.Password, version, user.Status, strings.Join(user.Roles, "\", \""))
}

// Write-only password configuration with password update
func testAccResourceSecurityUserWriteOnlyConfigWithPasswordUpdate(user security.User, firstName, newPassword string, version int) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid              = "%s"
	firstname           = "%s"
	lastname            = "%s"
	email               = "%s"
	password_wo         = "%s"
	password_wo_version = %d
	status              = "%s"
	roles               = ["%s"]
}
`, user.UserID, firstName, user.LastName, user.EmailAddress, newPassword, version, user.Status, strings.Join(user.Roles, "\", \""))
}

// Configuration that should cause conflict error
func testAccResourceSecurityUserConflictConfig(user security.User) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid              = "%s"
	firstname           = "%s"
	lastname            = "%s"
	email               = "%s"
	password            = "%s"
	password_wo         = "%s"
	password_wo_version = 1
	status              = "%s"
	roles               = ["%s"]
}
`, user.UserID, user.FirstName, user.LastName, user.EmailAddress, user.Password, user.Password, user.Status, strings.Join(user.Roles, "\", \""))
}

// Configuration with password_wo but no version (should fail)
func testAccResourceSecurityUserWriteOnlyNoVersionConfig(user security.User) string {
	return fmt.Sprintf(`
resource "nexus_security_user" "acceptance" {
	userid      = "%s"
	firstname   = "%s"
	lastname    = "%s"
	email       = "%s"
	password_wo = "%s"
	status      = "%s"
	roles       = ["%s"]
}
`, user.UserID, user.FirstName, user.LastName, user.EmailAddress, user.Password, user.Status, strings.Join(user.Roles, "\", \""))
}
