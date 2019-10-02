


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
    note_privileges_id integer PRIMARY KEY NOT NULL,
    note_id integer,
    user_name character varying(50),
    read CHAR(1),
    write CHAR(1)
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

INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Mohammad','Vaughn','Alec_Floyd','password'),('Curran','Cochran','Rose_Daugherty','password'),('Yoshio','Bernard','Raphael_Sexton','password'),('Hiram','Matthews','Kareem_Mosley','password'),('Kenyon','Wall','Graham_Mejia','password'),('Carson','Gillespie','Kelsey_Hutchinson','password'),('Hayes','Vinson','Todd_Fox','password'),('Jermaine','Alvarado','Serina_Rios','password'),('Andrew','House','Ralph_Blanchard','password'),('Rooney','Fowler','Allistair_Wiggins','password'),('Abbot','Greene','Porter_Long','password'),('Brandon','Terrell','Daniel_Massey','password'),('Donovan','Morris','Forrest_Stafford','password'),('Isaac','Gomez','Griffith_Dean','password'),('Plato','Myers','Griffith_Young','password'),('Omar','Caldwell','Nevada_Marsh','password'),('Aquila','Wyatt','Hall_Owen','password'),('Sean','Vincent','Cameran_Warner','password'),('Jonah','Rodriguez','Gannon_Cantrell','password'),('Zahir','Olsen','Audra_Summers','password');
INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Jameson','Huffman','Kasimir_Goodwin','password'),('Talon','Bernard','Stella_Nichols','password'),('Burton','Cruz','Kelly_Solomon','password'),('Arthur','Hopper','Natalie_Head','password'),('Armando','Sherman','Lee_Boyer','password'),('Finn','Reyes','Dillon_Myers','password'),('Cyrus','Gray','Jeanette_Britt','password'),('Tyrone','Wilkerson','Emmanuel_Perry','password'),('Elliott','Winters','Kevin_Donaldson','password'),('Sylvester','Nunez','Sylvester_Galloway','password'),('Aidan','Holt','Maggy_Ball','password'),('Malachi','Rodriguez','Theodore_Chambers','password'),('Cain','Poole','Silas_Nguyen','password'),('Keaton','Alvarado','Orson_Baldwin','password'),('Walter','Chapman','Fredericka_Mullen','password'),('Malachi','Mccarty','Amity_Byers','password'),('Cullen','Mclaughlin','Justin_Preston','password'),('Ferdinand','Cline','Allegra_Roach','password'),('Marsden','Woods','Devin_Neal','password'),('Gregory','Fitzgerald','Patience_Underwood','password');



