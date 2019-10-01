


CREATE TABLE _user (
  user_id integer PRIMARY KEY NOT NULL,
  given_name varchar(40),
  family_name varchar(40),
  user_name varchar(40),
  password varchar(40)
);


-- ------------------------------------------------

CREATE TABLE _note(
    note_id integer PRIMARY KEY NOT NULL,
    user_id integer,
    title VARCHAR(40),
    body VARCHAR(250),
    date_created DATE
);

-- ------------------------------------------------

CREATE TABLE _note_privileges(
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

INSERT INTO "_user" (user_id,given_name,family_name,user_name,password) VALUES (100,'Theodore','Contreras','Vance_Guerra','password'),(101,'Ishmael','Forbes','Lewis_Stanley','password'),(102,'Ezra','Hoover','Warren_Duke','password'),(103,'Hedley','Benson','Lane_Rich','password'),(104,'Yardley','Clark','Josiah_Slater','password'),(105,'Murphy','Boyer','Joel_Rocha','password'),(106,'Herman','Barber','Tana_Cabrera','password'),(107,'Hector','Vazquez','Sheila_Strong','password'),(108,'Callum','Lynn','Roanna_Saunders','password'),(109,'Ezekiel','Nelson','Kyle_Dudley','password'),(110,'Jermaine','Hardy','Randall_Kemp','password'),(111,'Henry','Ramirez','Barbara_Morales','password'),(112,'Zane','Reilly','David_Saunders','password'),(113,'Octavius','Terrell','Davis_Witt','password'),(114,'Jerry','Young','Quynn_Stanton','password'),(115,'Tyrone','Harmon','Isaac_Wilson','password'),(116,'Davis','Guerra','Randall_Hebert','password'),(117,'Grady','Williams','Amethyst_Olsen','password'),(118,'Phelan','Newman','Declan_May','password'),(119,'Tanek','Bridges','Nicole_Vargas','password');
INSERT INTO "_user" (user_id,given_name,family_name,user_name,password) VALUES (120,'Bevis','Webster','Dawn_Henry','password'),(121,'Dane','Berger','Jeanette_Hutchinson','password'),(122,'Zeus','Martin','Lamar_Madden','password'),(123,'Michael','Mercer','Ryder_Kelly','password'),(124,'Keefe','Hutchinson','Alexa_Gregory','password'),(125,'Driscoll','Pearson','Fritz_French','password'),(126,'Fritz','Stafford','Kyle_Savage','password'),(127,'Declan','Adams','Sylvia_Bradley','password'),(128,'Thaddeus','Hurley','Clio_Davenport','password'),(129,'Nissim','Luna','Colin_Riley','password'),(130,'Mason','Whitley','Ignatius_Crane','password'),(131,'Hamilton','Webster','Veronica_Hoffman','password'),(132,'George','Noel','Naomi_Palmer','password'),(133,'Plato','Nash','Eliana_Sutton','password'),(134,'Carson','Bond','Calvin_Hunt','password'),(135,'Felix','Church','Honorato_Landry','password'),(136,'Joshua','Richmond','Clare_Mckenzie','password'),(137,'Jackson','Blackburn','Dexter_Hurley','password'),(138,'Tarik','Alexander','Allegra_Peters','password'),(139,'Neil','Goff','Cade_Sullivan','password');
INSERT INTO "_user" (user_id,given_name,family_name,user_name,password) VALUES (140,'Myles','Ratliff','Lucian_Oconnor','password'),(141,'Isaiah','Hayden','Chadwick_Doyle','password'),(142,'Isaac','Terry','Victoria_Nash','password'),(143,'Jordan','Kennedy','Dieter_Dodson','password'),(144,'Acton','Hodges','Joelle_Nielsen','password'),(145,'Dante','Campbell','Gary_Odom','password'),(146,'Simon','Maynard','Sheila_Clay','password'),(147,'Tobias','Taylor','Ocean_Mcneil','password'),(148,'Gavin','Spence','Clarke_Goff','password'),(149,'Carlos','Walton','Giacomo_Melendez','password'),(150,'Carl','Travis','Kendall_Hyde','password'),(151,'Preston','Atkins','Lesley_Mckenzie','password'),(152,'Vance','Fuller','Alyssa_Tucker','password'),(153,'Leo','Potts','Chaim_Dalton','password'),(154,'Lawrence','Fulton','Kiara_Paul','password'),(155,'Channing','Mccullough','Lillian_Michael','password'),(156,'Aristotle','Shaw','Dawn_Clayton','password'),(157,'Giacomo','Rivas','Elliott_Lynch','password'),(158,'Tiger','Shepherd','Daniel_Brooks','password'),(159,'Blaze','Burton','Christian_Bender','password');
INSERT INTO "_user" (user_id,given_name,family_name,user_name,password) VALUES (160,'Brendan','Kinney','Oscar_Underwood','password'),(161,'Giacomo','Miller','Heidi_Caldwell','password'),(162,'Harper','Flores','Elton_Morgan','password'),(163,'Quentin','Roach','Maya_Albert','password'),(164,'Dexter','Joyce','Vincent_Mcguire','password'),(165,'Ashton','Crosby','Lyle_Ayers','password'),(166,'Kieran','Wilson','Joshua_Flynn','password'),(167,'Camden','Carson','Signe_Mejia','password'),(168,'Baker','Eaton','Dustin_Valdez','password'),(169,'Elijah','Diaz','Constance_Solis','password'),(170,'Cameron','Hutchinson','Stewart_Estes','password'),(171,'Drake','Rhodes','Alden_Hernandez','password'),(172,'Rashad','Bowen','Hiroko_Watson','password'),(173,'Samuel','Serrano','Maryam_Diaz','password'),(174,'Jin','Chan','Lynn_Carver','password'),(175,'Jordan','Fitzgerald','Natalie_Cantrell','password'),(176,'Linus','Lester','Belle_Terry','password'),(177,'Moses','Floyd','Constance_Wilkins','password'),(178,'Marsden','Powell','Nicholas_Dickson','password'),(179,'Thomas','Slater','Cooper_Huber','password');
INSERT INTO "_user" (user_id,given_name,family_name,user_name,password) VALUES (180,'Merritt','Randall','Lavinia_Hodge','password'),(181,'Moses','Wilder','Leonard_Rowe','password'),(182,'Simon','Mcneil','Ulysses_Wright','password'),(183,'Forrest','Erickson','Kelsey_Abbott','password'),(184,'Neil','Hartman','Myles_Washington','password'),(185,'Laith','Holt','Hannah_Lowery','password'),(186,'Ivor','Lara','Christine_Carpenter','password'),(187,'Lucius','Ballard','Norman_Griffin','password'),(188,'Silas','Smith','Brady_Hammond','password'),(189,'Lars','Kelley','Fritz_Dixon','password'),(190,'Patrick','Mayo','Malcolm_Harrison','password'),(191,'Daniel','Tucker','Moana_Blevins','password'),(192,'Hunter','Bullock','Andrew_Kim','password'),(193,'Abel','Durham','Nathaniel_Rose','password'),(194,'Colby','Santiago','Lenore_Mcguire','password'),(195,'Odysseus','Glover','Gisela_Russell','password'),(196,'Maxwell','Hardin','Nasim_Holman','password'),(197,'Caesar','Golden','George_Vargas','password'),(198,'Sebastian','Stevenson','Alec_Koch','password'),(199,'Galvin','Bolton','Connor_Calhoun','password');

