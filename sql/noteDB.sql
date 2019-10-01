


CREATE TABLE _user ( -- added the underscore so we dont have to always add quotes in code
  given_name varchar(40),
  family_name varchar(40),
  userName  character varying(50) NOT NULL PRIMARY KEY,
  password varchar(40)
);


-- ------------------------------------------------

CREATE TABLE _note( -- added underscore here to keep naming convention
    note_id integer PRIMARY KEY NOT NULL,
    user_id integer,
    title VARCHAR(40),
    body VARCHAR(250),
    date_created DATE
);

-- ------------------------------------------------

CREATE TABLE _note_privileges( -- added underscore here to keep naming convention
    note_privileges_id integer PRIMARY KEY NOT NULL,
    note_id integer,
    user_id integer,
    read CHAR(1),
    write CHAR(1)
);

-- ------------------------------------------------

ALTER TABLE _note ADD  
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES _user (user_id);

-- ------------------------------------------------

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES _user (user_id);

ALTER TABLE "_note_privileges" ADD  
    CONSTRAINT note_id FOREIGN KEY (note_id)
        REFERENCES "_note" (note_id);


-- ------------------------------------------------

INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Mohammad','Vaughn','Alec_Floyd','password'),('Curran','Cochran','Rose_Daugherty','password'),('Yoshio','Bernard','Raphael_Sexton','password'),('Hiram','Matthews','Kareem_Mosley','password'),('Kenyon','Wall','Graham_Mejia','password'),('Carson','Gillespie','Kelsey_Hutchinson','password'),('Hayes','Vinson','Todd_Fox','password'),('Jermaine','Alvarado','Serina_Rios','password'),('Andrew','House','Ralph_Blanchard','password'),('Rooney','Fowler','Allistair_Wiggins','password'),('Abbot','Greene','Porter_Long','password'),('Brandon','Terrell','Daniel_Massey','password'),('Donovan','Morris','Forrest_Stafford','password'),('Isaac','Gomez','Griffith_Dean','password'),('Plato','Myers','Griffith_Young','password'),('Omar','Caldwell','Nevada_Marsh','password'),('Aquila','Wyatt','Hall_Owen','password'),('Sean','Vincent','Cameran_Warner','password'),('Jonah','Rodriguez','Gannon_Cantrell','password'),('Zahir','Olsen','Audra_Summers','password');
INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Jameson','Huffman','Kasimir_Goodwin','password'),('Talon','Bernard','Stella_Nichols','password'),('Burton','Cruz','Kelly_Solomon','password'),('Arthur','Hopper','Natalie_Head','password'),('Armando','Sherman','Lee_Boyer','password'),('Finn','Reyes','Dillon_Myers','password'),('Cyrus','Gray','Jeanette_Britt','password'),('Tyrone','Wilkerson','Emmanuel_Perry','password'),('Elliott','Winters','Kevin_Donaldson','password'),('Sylvester','Nunez','Sylvester_Galloway','password'),('Aidan','Holt','Maggy_Ball','password'),('Malachi','Rodriguez','Theodore_Chambers','password'),('Cain','Poole','Silas_Nguyen','password'),('Keaton','Alvarado','Orson_Baldwin','password'),('Walter','Chapman','Fredericka_Mullen','password'),('Malachi','Mccarty','Amity_Byers','password'),('Cullen','Mclaughlin','Justin_Preston','password'),('Ferdinand','Cline','Allegra_Roach','password'),('Marsden','Woods','Devin_Neal','password'),('Gregory','Fitzgerald','Patience_Underwood','password');
INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Carter','Wyatt','Thaddeus_Allison','password'),('Ali','Keith','Haley_Flynn','password'),('Fritz','Kane','Nigel_Rosario','password'),('Kareem','Nicholson','Garth_Drake','password'),('Benedict','Chan','Ashely_Moody','password'),('Jacob','Moran','Lareina_Carey','password'),('Eric','Mckenzie','Nissim_Bradford','password'),('Jelani','Osborne','Keegan_Owens','password'),('Amir','Dean','Lester_Burns','password'),('Silas','Saunders','Marvin_Jensen','password'),('Drew','Hinton','Chastity_Dale','password'),('Jerry','Ferrell','Farrah_Church','password'),('Kevin','Sexton','Piper_Guerrero','password'),('Wylie','Bartlett','Clark_Chan','password'),('Abraham','Garza','Yvonne_Peck','password'),('Xander','Holman','Jocelyn_Bird','password'),('Caleb','Castaneda','Melissa_Hopper','password'),('Richard','Schroeder','Wynne_Middleton','password'),('Giacomo','Simon','Kylee_Berger','password'),('Melvin','Carroll','Odysseus_Landry','password');
INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Kieran','Maddox','Cally_Mcintosh','password'),('Giacomo','Howard','Francesca_Stephenson','password'),('Addison','Whitfield','Salvador_Vaughn','password'),('Colby','Sanford','Aimee_Rojas','password'),('Dane','Osborn','Andrew_Stevens','password'),('Noah','Vaughn','Isaac_Sparks','password'),('Honorato','Price','Steven_Goodwin','password'),('Gannon','Douglas','David_Joyner','password'),('Geoffrey','Carr','Desiree_Mosley','password'),('Bert','Webb','Lydia_Joyce','password'),('Carl','Suarez','Stewart_Wilkerson','password'),('Devin','Buckner','Hedda_Warner','password'),('Adrian','Flowers','Thane_Morin','password'),('Bruce','Kerr','Quon_Roach','password'),('Barclay','Rasmussen','Melissa_Guerra','password'),('Andrew','Blake','Chava_Luna','password'),('Nissim','Phillips','Lisandra_Cortez','password'),('Brett','Pacheco','Xantha_Pratt','password'),('Seth','Calhoun','Price_Rowe','password'),('Chandler','Whitley','Olivia_Solomon','password');
INSERT INTO "_user" (given_name,family_name,user_name,password) VALUES ('Dustin','Pitts','Simone_Meadows','password'),('Ethan','Black','Martena_Faulkner','password'),('Logan','Patel','Alisa_Logan','password'),('Quamar','Barry','Hilary_Rhodes','password'),('Garrison','Mullins','Imelda_Oneill','password'),('Igor','Barber','Colorado_Allison','password'),('Holmes','Combs','Sharon_House','password'),('Forrest','Carroll','Ulla_Hendricks','password'),('Tucker','Harmon','Thane_Wyatt','password'),('Kenneth','Vega','Stuart_Ferguson','password'),('John','Garrett','Melanie_Smith','password'),('Forrest','Obrien','Fay_Porter','password'),('Ryan','Buckley','Acton_Gutierrez','password'),('Norman','Gallagher','Sybil_Carpenter','password'),('Kelly','Miranda','Mason_Rodriquez','password'),('Reuben','Schmidt','Cody_Lewis','password'),('Timon','Emerson','Zachary_Howe','password'),('Edan','Grant','Jordan_Byrd','password'),('Nathaniel','Rodgers','Todd_Calderon','password'),('Josiah','Camacho','Selma_Delgado','password');
