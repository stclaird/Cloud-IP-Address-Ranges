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
	578553856,
	578554879,
	"https://www.gstatic.com/ipranges/cloud.json",
	"google",
	"IPv4"
);