create database hami character set utf8mb4;


create table resource_pool(
                              id          bigint primary key not null auto_increment,
                              pool_name   varchar(128)       not null,
                              create_time timestamp default current_timestamp,
                              update_time  timestamp default current_timestamp on update current_timestamp
);

create table nodes(
                      id bigint primary key not null auto_increment,
                      pool_id bigint not null,
                      node_name varchar(128) not null,
                      node_ip varchar(32) not null,
                      create_time timestamp default current_timestamp,
                      update_time timestamp default current_timestamp on update current_timestamp
);

INSERT INTO hami.resource_pool (pool_name) VALUES ('大模型资源池');
