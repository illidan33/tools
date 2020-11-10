package clients

//go:generate tools gen client --url "http://192.168.1.116:8080/swagger/swagger/doc.json" -n "Fpx"

//go:generate tools gen client --url "http://192.168.1.175:8080/swagger/swagger/doc.json" -n "Gkspg"

//go:generate tools gen client --url "http://47.254.245.66:10883/swagger/swagger/doc.json" -n "Gkspg"

//go:generate tools gen client --url "./doc/gkspg.json" -n "Gkspg" -d true
