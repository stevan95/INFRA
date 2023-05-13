CREATE DATABASE ec2informations;

CREATE TABLE instances (
   InstanceID varchar(50),
   PrivateIpAddress varchar(50),
   PrivateDnsName varchar(50),
   SubnetId varchar(50),
   PublicIpAddress varchar(50),
   InstanceType varchar(50),
   PRIMARY KEY(InstanceID)
);

CREATE TABLE instances_tags (
   TagID INT GENERATED ALWAYS AS IDENTITY,
   InstanceID varchar(50) REFERENCES instances (InstanceID) ON DELETE CASCADE,
   TagsKey varchar(50),
   TagsName varchar(50),
   PRIMARY KEY (TagID)
);
