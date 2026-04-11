--Database
CREATE TABLE IF NOT EXISTS net (
 	net_id INTEGER PRIMARY KEY,
 	net TEXT NOT NULL,
 	start_ip INT NOT NULL,
 	end_ip INT NOT NULL,
 	url TEXT NOT NULL,
 	cloudplatform TEXT NOT NULL,
 	iptype TEXT NOT NULL
 	);

INSERT INTO net (net_id, net, start_ip, end_ip, url, cloudplatform, iptype)
VALUES (
	NULL, 
	"34.124.8.0/22",
	579665920,
	579666943,
	"https://www.gstatic.com/ipranges/cloud.json",
	"google",
	"IPv4"
);

-- DuckDB example (use TEXT for start/end to support IPv6)
CREATE TABLE IF NOT EXISTS net_duckdb (
	net TEXT PRIMARY KEY,
	start_ip TEXT NOT NULL,
	end_ip TEXT NOT NULL,
	url TEXT NOT NULL,
	cloudplatform TEXT NOT NULL,
	iptype TEXT NOT NULL
	);

INSERT INTO net_duckdb (net, start_ip, end_ip, url, cloudplatform, iptype)
VALUES (
	"2001:db8::/32",
	"2001:db8::",
	"2001:db8:ffff:ffff:ffff:ffff:ffff",
	"https://example.test",
	"example",
	"IPv6"
);