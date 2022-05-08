create table parent(parent_id INT PRIMARY KEY, name VARCHAR(20), current_baby INT);

create table baby(baby_id SERIAL PRIMARY KEY, parent_id INT REFERENCES parent, name VARCHAR (20), birth DATE);

create table sleep(id SERIAL PRIMARY KEY, baby_id INT REFERENCES baby, start TIMESTAMP NOT NULL, sleep_end TIMESTAMP);

create table eat(id SERIAL PRIMARY KEY, baby_id INT REFERENCES baby, start TIMESTAMP NOT NULL, description VARCHAR);

create table pampers(id SERIAL PRIMARY KEY, baby_id INT REFERENCES baby, start TIMESTAMP NOT NULL, state INT NOT NULL);
