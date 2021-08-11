/*
 * Author: Luis Guillermo Gómez Galeano
 *
 * Initial DML for the data base.
 */

INSERT INTO users(id, first_name, last_name, email, password)
VALUES (1, 'Daniel', 'Valencia', 'daniel@gmail.com',
        '$argon2id$v=19$m=102400,t=2,p=8$PI0LWDyTIMtJKz34rqX1sw$3heILJEia4RQSgNlIHlFOQ'),
       (2, 'Carlos', 'Roldán', 'carlos@gmail.com',
        '$argon2id$v=19$m=102400,t=2,p=8$pmgpBYfEswnpFO90j5P8aw$nuPlvspaOM6pGTG8dVd9ww'),
       (3, 'Ana', 'Guerrero', 'ana@gmail.com',
        '$argon2id$v=19$m=102400,t=2,p=8$w2EvBoIyqX6Tnp00WqHiSQ$ooBkRg9imp/YzCRi/s0f7A');

SELECT setval('users_id_seq', 3);
