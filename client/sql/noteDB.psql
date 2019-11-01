

DROP TABLE IF EXISTS _user CASCADE;
CREATE TABLE _user( -- added the underscore so we dont have to always add quotes in code
  given_name varchar(40),
  family_name varchar(40),
  user_name  character varying(50) NOT NULL PRIMARY KEY,
  password varchar(40),
  email varchar(40),
  session_id varchar(100)
);


-- ------------------------------------------------

DROP TABLE IF EXISTS _note CASCADE;
CREATE TABLE _note( -- added underscore here to keep naming convention
    note_id serial PRIMARY KEY NOT NULL,
    note_owner character varying(50),
    title VARCHAR(40),
    body VARCHAR(250),
    date_created DATE
);

-- ------------------------------------------------

DROP TABLE IF EXISTS _note_privileges CASCADE;
CREATE TABLE _note_privileges( -- added underscore here to keep naming convention
    note_privileges_id serial PRIMARY KEY NOT NULL,
    note_id integer,
    user_name character varying(50),
    read CHAR(1), -- t for true  f for false
    write CHAR(1) -- t for true  f for false
);

DROP TABLE IF EXISTS _group CASCADE;
CREATE TABLE _group(
    group_id serial PRIMARY KEY NOT NULL,
    group_title VARCHAR(40),
    read CHAR(1),
    write CHAR(1),
    group_owner VARCHAR(40)
);

DROP TABLE IF EXISTS _group_user CASCADE;
CREATE TABLE _group_user(
    group_user_id serial PRIMARY KEY NOT NULL,
    group_id integer,
    user_name VARCHAR(40)
);

-- ------------------------------------------------

ALTER TABLE _note ADD  
    CONSTRAINT note_owner FOREIGN KEY (note_owner)
        REFERENCES _user (user_name);

-- ------------------------------------------------

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT user_name FOREIGN KEY (user_name)
        REFERENCES _user (user_name);

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT note_id FOREIGN KEY (note_id)
        REFERENCES "_note" (note_id);


-- ------------------------------------------------

ALTER TABLE _group ADD  
    CONSTRAINT group_owner FOREIGN KEY (group_owner)
        REFERENCES _user (user_name);


-- ------------------------------------------------

ALTER TABLE "_group_user" ADD  
    CONSTRAINT user_name FOREIGN KEY (user_name)
        REFERENCES _user (user_name);

ALTER TABLE "_group_user" ADD  
    CONSTRAINT group_id FOREIGN KEY (group_id)
        REFERENCES "_group" (group_id);


-- ------------------------------------------------


INSERT INTO "_user" (given_name,family_name,user_name,password,email) 
VALUES 
('Travis', 'Falkenberg', 'Trav3', '1234', 'travis.falkenberg141@gmail.com'),
('Mohammad','Vaughn','Vaughn1','password', 'fakeemail1@gmai2.com'),
('Curran','Cochran','Curran85','password', 'fakeemail1@gmai3.com'),
('Yoshio','Bernard','BernardsProfile','password', 'fakeemail1@gmai4.com'),
('Hiram','Matthews','Matthews_Hiram','password', 'fakeemail1@gmai5.com'),
('Kenyon','Wall','Grand_Kenyon','password', 'fakeemail1@gmai6.com'),
('Carson','Gillespie','Carson123','password', 'fakeemail1@gmai7.com'),
('Hayes','Vinson','Hayes45','password', 'fakeemail1@gmai8.com'),
('Jermaine','Alvarado','Jermaine321','password', 'fakeemail1@gmai9.com'),
('Andrew','House','Andrew222','password', 'fakeemail1@gmail0.com'),
('Rooney','Fowler','Wayne_Rooney','password', 'fakeemail1@gmail1.com'),
('Abbot','Greene','Abbot_Time','password', 'fakeemail1@gmail2.com'),
('Brandon','Terrell','Daniel_Massey','password', 'fakeemail1@gmail3.com'),
('Donovan','Morris','Forrest_Stafford','password', 'fakeemail1@gmail4.com'),
('Isaac','Gomez','Griffith_Dean','password', 'fakeemail1@gmail5.com'),
('Plato','Myers','Griffith_Young','password', 'fakeemail1@gmail6.com'),
('Omar','Caldwell','Nevada_Marsh','password', 'fakeemail1@gmail7.com'),
('Aquila','Wyatt','Hall_Owen','password', 'fakeemail1@gmail8.com'),
('Sean','Vincent','Cameran_Warner','password', 'fakeemail1@gmail9.com'),
('Jonah','Rodriguez','Gannon_Cantrell','password', 'fakeemail1@gmai20.com'),
('Zahir','Olsen','Audra_Summers','password', 'fakeemail1@gmai21.com');



INSERT INTO _note (note_owner, title, body, date_created)
VALUES
('Vaughn1', 'note test', 'A note i created to test notes, this is nothing interesting', date('now')),
('Vaughn1', 'Vaughn note', 'Vaughn wrote this note, he has added weird words like twist or hippo', date('now')),
('Vaughn1', 'note i worte', 'Vaughn note with the word twist', date('now')),
('Trav3', 'NoteBookApp To do List', 'This is a list of things we need to do for the webapp', date('now'));


INSERT INTO _note_privileges(note_id, user_name, read, write)
VALUES
(1, 'Trav3', 't', 't');

-- for word match search 

-- SELECT * FROM _note
-- WHERE body ~ '\yhippo\y';

-- SELECT * FROM _note
-- WHERE body ~ '\ytwist\y';

-- partial word match search

-- SELECT * FROM _note
-- WHERE body ~ 'twi:*';

-- SELECT * FROM _note
-- WHERE body ~ 'ppo:*';