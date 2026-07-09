resource "nexus_security_oidc" "example" {
  client_id         = "nexus"
  client_secret     = "very-secret"
  authorization_url = "https://idp.example.test/oauth2/authorize"
  token_url         = "https://idp.example.test/oauth2/token"
  jwks_url          = "https://idp.example.test/.well-known/jwks.json"
  jws_algorithm     = "RS256"
  username_claim    = "preferred_username"
  groups_claim      = "groups"
  logout_url        = "https://idp.example.test/oauth2/logout"
  first_name_claim  = "given_name"
  last_name_claim   = "family_name"
  email_claim       = "email"
  use_trust_store   = false
}
