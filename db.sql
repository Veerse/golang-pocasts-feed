drop table episodes;
drop table podcasts;

create table podcasts (
                          id serial primary key,
                          title varchar,
                          description varchar,
                          image varchar,
                          language varchar,
                          category varchar,
                          author_name varchar,
                          author_email varchar,
                          link varchar,
                          owner varchar
);


insert into podcasts (title, description, image, language, category, author_name, author_email, link, owner)
VALUES ('La voix de l''homme trouble', 'Dans ce podcast l''homme trouble s''exprime et nous raconte des trucs dont on se bats les roubistoles',
        'http://yapas.jpg/', 'fr-FR', 'mmh', 'the illusive man', 'nab@fakemail.net', 'http://mylink.net/', 'muslimy');

insert into podcasts (title, description, image, language, category, author_name, author_email, link, owner)
VALUES ('Le destin sombre de nassim le malefique', 'Ce podcast c''est pas pour les pieds tendres attention',
        'http://aze.jpg/', 'fr-FR', 'euh', 'the illusive man', 'nab@fakemail.net', 'http://thelink.net/', 'muslimy');

create table episodes (
                          id serial primary key,
                          podcast_id int,
                          title varchar,
                          url varchar,
                          length int,
                          type varchar,
                          guid varchar,
                          pub_date date,
                          description varchar,
                          episode_url varchar,
                          image varchar,

                          foreign key (podcast_id) references podcasts(id)
);

insert into episodes (podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image) VALUES
(1, 'la naissance', 'http://ghub.com/', 1224, 'audio/mp3', 1, 'Wed, 15 Jun 2019 19:00:00 GMT', 'l''homme trouble explique sa ptn de naissance', 'http://yapas.com', 'http://noplus.png');

insert into episodes (podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image) VALUES
(1, 'l''enfacnce', 'https://link2.net/', 2446, 'audio/mp3', 2, 'Wed, 22 Jun 2019 19:00:00 GMT', 'dans cette episode il nous raconte son enface difficile dans les rues de chicago', 'http://lien.mp3/', 'http://lefauxlien.jpg/');

insert into episodes (podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image) VALUES
(2, 'le jour où il a flingue 14 bougzer vergogneless', 'http://nope.com/', 4403, 'audio/mp3', 3, 'Wed, 16 Jun 2019 19:00:00 GMT', 'La vergogne etait totallement absente ce jour la je vous le jure', 'http://lelien.mp3/', 'http://quedalle.jpg/');