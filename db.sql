create table podcasts (
                          id serial primary key,
                          title varchar,
                          description varchar,
                          image varchar,
                          language varchar,
                          category varchar,
                          author varchar,
                          link varchar,
                          owner varchar
);


insert into podcasts (title, description, image, language, category, author, link, owner)
VALUES ('La voix de l''homme trouble', 'Dans ce podcast l''homme trouble s''exprime et nous raconte des trucs dont on se bats les roubistoles',
        'y en a pas', 'fr-FR', 'mmh', 'the illusive man', 'tjr pas', 'muslimy');

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

                          foreign key (podcast_id) references podcast(id)
);

insert into episodes (podcast_id, title, url, length, type, guid, pub_date, description, episode_url, image) VALUES
(1, 'la naissance', 'aucun', 1224, 'audio/mp3', 1, 'Wed, 15 Jun 2019 19:00:00 GMT', 'l''homme trouble explique sa ptn de naissance', 'y a pas', 'non plus')