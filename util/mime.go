package util

type MIME string

const (
    MIMEJson              MIME = "application/json"
    MIMEPostForm          MIME = "application/x-www-form-urlencoded"
    MIMEMultipartPOSTForm MIME = "multipart/form-data"
)
