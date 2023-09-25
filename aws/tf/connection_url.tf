resource "aws_secretsmanager_secret" "db_admin_pg" {
  name = "${var.name}/conn_url"
  tags = var.tags
}

resource "aws_secretsmanager_secret_version" "db_admin_pg" {
  secret_id     = aws_secretsmanager_secret.db_admin_pg.id
  port_substring = endswith(var.protocol, "srv") ? "" : ":${var.port}"
  secret_string = "${var.protocol}://${urlencode(var.username)}:${urlencode(var.password)}@${var.host}/${urlencode(var.database)}"
}
