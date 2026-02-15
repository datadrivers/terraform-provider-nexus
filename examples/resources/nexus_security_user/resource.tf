resource "nexus_security_user" "admin" {
  userid    = "admin"
  firstname = "Administrator"
  lastname  = "User"
  email     = "nexus@example.com"
  password  = "admin123"
  roles     = ["nx-admin"]
  status    = "active"
  source    = "default"
}

resource "nexus_security_user" "user_password_wo" {
  userid              = "user_password_wo"
  firstname           = "Administrator"
  lastname            = "User"
  email               = "nexus1@example.com"
  password_wo         = "admin123"   # This password value don't save to state
  password_wo_version = 1            # Incriment version, for update password  
  roles               = ["nx-admin"]
  status              = "active"
}

ephemeral "random_password" "password" {
  length           = 16
}

resource "nexus_security_user" "user_password_from_ephemeral" {
  userid              = "user_ephemeral_password"
  firstname           = "ephemeral"
  lastname            = "User"
  email               = "ephemeral@example.com"
  password_wo         = ephemeral.random_password.password  # Use ephemeral value
  password_wo_version = 1                                   # Incriment version, for update password  
  roles               = ["nx-admin"]
  status              = "active"
}
