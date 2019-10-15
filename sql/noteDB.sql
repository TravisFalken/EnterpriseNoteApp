


CREATE TABLE _user ( -- added the underscore so we dont have to always add quotes in code
  given_name varchar(40),
  family_name varchar(40),
  user_name  character varying(50) NOT NULL PRIMARY KEY,
  password varchar(40)
);


-- ------------------------------------------------

CREATE TABLE _note( -- added underscore here to keep naming convention
    note_id integer PRIMARY KEY NOT NULL,
    user_name character varying(50),
    title VARCHAR(40),
    body VARCHAR(250),
    date_created DATE
);

-- ------------------------------------------------

CREATE TABLE _note_privileges( -- added underscore here to keep naming convention
    note_privileges_id serial PRIMARY KEY NOT NULL,
    note_id integer,
    user_name character varying(50),
    read CHAR(1), -- t for true  f for false
    write CHAR(1) -- t for true  f for false
);

-- ------------------------------------------------

ALTER TABLE _note ADD  
    CONSTRAINT user_name FOREIGN KEY (user_name)
        REFERENCES _user (user_name);

-- ------------------------------------------------

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT user_name FOREIGN KEY (user_name)
        REFERENCES _user (user_name);

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT note_id FOREIGN KEY (note_id)
        REFERENCES "_note" (note_id);


-- ------------------------------------------------

INSERT INTO "_user" (given_name,family_name,user_name,password) 
VALUES 
('Mohammad','Vaughn','Vaughn1','password'),
('Curran','Cochran','Curran85','password'),
('Yoshio','Bernard','BernardsProfile','password'),
('Hiram','Matthews','Matthews_Hiram','password'),
('Kenyon','Wall','Grand_Kenyon','password'),
('Carson','Gillespie','Carson123','password'),
('Hayes','Vinson','Hayes45','password'),
('Jermaine','Alvarado','Jermaine321','password'),
('Andrew','House','Andrew222','password'),
('Rooney','Fowler','Wayne_Rooney','password'),
('Abbot','Greene','Abbot_Time','password'),
('Brandon','Terrell','Daniel_Massey','password'),
('Donovan','Morris','Forrest_Stafford','password'),
('Isaac','Gomez','Griffith_Dean','password'),
('Plato','Myers','Griffith_Young','password'),
('Omar','Caldwell','Nevada_Marsh','password'),
('Aquila','Wyatt','Hall_Owen','password'),
('Sean','Vincent','Cameran_Warner','password'),
('Jonah','Rodriguez','Gannon_Cantrell','password'),
('Zahir','Olsen','Audra_Summers','password');



INSERT INTO _note (note_id, user_name, title, body, date_created)
VALUES
(100, 'Vaughn1', 'note test', 'A note i created to test notes, this is nothing interesting', date('now')),
(101, 'Curran85', 'currans note', 'curran wrothe this note, he has added weird words like twist or hippo', date('now')),
(102, 'Grand_Kenyon', 'note i worte', 'kenyons note with the word twist', date('now'));

-- for word match search 

SELECT * FROM _note
WHERE body ~ '\yhippo\y';

SELECT * FROM _note
WHERE body ~ '\ytwist\y';

-- partial word match search

SELECT * FROM _note
WHERE body ~ 'twi:*';

SELECT * FROM _note
WHERE body ~ 'ppo:*';