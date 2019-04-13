# hs-micro-back

## Overview

This is a educational purpose app: a simple bloging like platform.

## Architecture

```
             ---------
             | EMAIL |
             ---------
                 ^
                 |
---------     --------     --------     -----------
| FRONT | --> | NATS | --> | BACK | --> | MariaDB |
---------     --------     --------     -----------
```

 - Front: a go frontend (gorilla, html/templatesn go-nats)
 - Back: a go backend (go-nats, database/sql)
 - Email: a python notification service

## Run

### Prep

Before you begin, prepare your database:

```
CREATE TABLE `post`(
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`post_title` varchar(64) DEFAULT NULL,
	`post_content` mediumtext,
	`post_date` timestamp NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
	) ENGINE=InnoDB DEFAULT CHARSET=latin1;
```

### Binary

```
$ go build -o app
$ export NATSURL="your_nats_url"               // default demo.nats.io
$ export NATSPORT="your_nats_port"             // default :4222
$ export NATSPOST="your_nats_post_channel"     // the channel used for posts, default zjnO12CgNkHD0IsuGd89zA
$ export NATSGET="your_nats_get_posts_channel" // the channel used get posts, default OWM7pKQNbXd7l75l21kOzA
$ export DBUSER="your_db_user"                 // default user
$ export DBPASS="your_db_password"             // default password
$ export DBHOST="your_db_host"                 // default 127.0.0.1
$ export DBPORT="your_db_port"                 // default :3306
$ export DBBASE="your_db_database"             // default blowofmouth
$ ./app
```

### Docker

```
$ export NATSURL="your_nats_url"               // default demo.nats.io
$ export NATSPORT="your_nats_port"             // default :4222
$ export NATSPOST="your_nats_post_channel"     // the channel used for posts, default zjnO12CgNkHD0IsuGd89zA
$ export NATSGET="your_nats_get_posts_channel" // the channel used get posts, default OWM7pKQNbXd7l75l21kOzA
$ export DBUSER="your_db_user"                 // default user
$ export DBPASS="your_db_password"             // default password
$ export DBHOST="your_db_host"                 // default 127.0.0.1
$ export DBPORT="your_db_port"                 // default :3306
$ export DBBASE="your_db_database"             // default blowofmouth
$ docker run -d -e NATSURL=${NATSURL} -e NATSPORT=${NATSPORT} -e NATSPOST=${NATSPOST} -e NATSGET=${NATSGET} \ 
    -e DBUSER=${DBUSER} -e DBPASS=${DBPASS} -e DBHOST=${DBHOST} -e DBPORT=${DBPORT} -e DBBASE=${DBBASE} \ 
    -p 8080:8080 jblaskovich/hs-micro-back:$release
```