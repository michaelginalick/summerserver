-- -- -- psql -d cs6400_fa18_team047 -f seed.sql
/*
DB seeding for corkboardit
Functions make things DRYer
*/

/*
Users
input = [[user_name, email, pin], ...]
*/

CREATE OR REPLACE FUNCTION create_users(text[][]) RETURNS void AS $$
DECLARE
  _user text[];
BEGIN
  FOREACH _user SLICE 1 IN ARRAY $1
  LOOP
    INSERT INTO "user" (user_name, email, pin) VALUES (_user[1], _user[2], _user[3]::int);
  END LOOP;
END;
$$ LANGUAGE plpgsql;


SELECT 'Creating Users...' FROM create_users(
    ARRAY[
         ['Ian Knox', 'ian3@gt.edu', '1234'],
         ['Michael Ginalick', 'michael@gt.edu', '1234'],
         ['Joe Blow', 'jb@123.com', '1234']
    ]
);


/*
Follows
input = [follower_email, followee_email]
*/

CREATE OR REPLACE FUNCTION create_follows(text[][]) RETURNS void AS $$
DECLARE
  _follows text[];
BEGIN
  FOREACH _follows SLICE 1 IN ARRAY $1
	LOOP
		INSERT INTO follow (
		    user_id,
		    followee_id
		)
		VALUES (
			(SELECT id AS _follower FROM "user" WHERE email = _follows[1]),
			(SELECT id AS _followee FROM "user" WHERE email = _follows[2])
		);
  END LOOP;
END;
$$ LANGUAGE plpgsql;


SELECT 'Creating Follows...' FROM create_follows(
    ARRAY[
            ['ian3@gt.edu','michael@gt.edu'],
            ['michael@gt.edu', 'jb@123.com'],
            ['jb@123.com', 'ian3@gt.edu']
    ]
);

/*
Categories - no need for a function
*/
SELECT 'Creating Categories...';
INSERT INTO category (category_name) VALUES
('Education'),
('People'),
('Sports'),
('Other'),
('Architecture'),
('Travel'),
('Pets'),
('Food & Drink'),
('Home & Garden'),
('Photography'),
('Technology'),
('Art');

/*
Public Corkboards
input = [[title, user_id], etc...]
Category is selected randomly
*/
CREATE OR REPLACE FUNCTION create_public_corkboards(text[][]) RETURNS void AS $$
DECLARE
  _corkboard text[];
BEGIN
  FOREACH _corkboard SLICE 1 IN ARRAY $1
	LOOP
	    WITH lookup_insert AS (
	        INSERT INTO corkboard(id)
	        VALUES (DEFAULT)
	        RETURNING id)
		INSERT INTO public_corkboard (
		    id,
            title,
            user_id,
            category_name)
		VALUES (
			(SELECT id FROM lookup_insert),
		    _corkboard[1],
			(SELECT id FROM "user" WHERE email = _corkboard[2]),
			(SELECT category_name FROM "category" ORDER BY RANDOM() LIMIT 1)
        );
  END LOOP;
END;
$$ LANGUAGE plpgsql;


SELECT 'Creating Public Corkboards...' FROM create_public_corkboards(
    ARRAY[
         ['Tranquility Base Hotel & Casino', 'ian3@gt.edu'],
         ['I Like It When You Sleep', 'michael@gt.edu'],
         ['Mellon Collie And The Infinite Sadness', 'jb@123.com'],
         ['Different Gear, Still Speeding', 'ian3@gt.edu'],
         ['The Hour Of Bewilderbeast', 'michael@gt.edu'],
         ['Icky Thump', 'jb@123.com'],
         ['Dreamboat Annie', 'ian3@gt.edu']
    ]
);


/*
Private Corkboards
input = [[title, user_id], etc...]
Category is selected randomly, password is always p@55w0rd
*/
CREATE OR REPLACE FUNCTION create_private_corkboards(text[][]) RETURNS void AS $$
DECLARE
  _corkboard text[];
BEGIN
  FOREACH _corkboard SLICE 1 IN ARRAY $1
	LOOP
	    WITH lookup_insert AS (
	        INSERT INTO corkboard(id)
	        VALUES (DEFAULT)
	        RETURNING id)
		INSERT INTO private_corkboard (
		    id,
            title,
            user_id,
            category_name,
            password)
		VALUES (
			(SELECT id FROM lookup_insert),
		    _corkboard[1],
			(SELECT id FROM "user" WHERE email = _corkboard[2]),
			(SELECT category_name FROM "category" ORDER BY RANDOM() LIMIT 1),
			'p@55w0rd'
        );
  END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT 'Creating Private Corkboards...' FROM create_private_corkboards(
    ARRAY[
         ['BMW 328i', 'jb@123.com'],
         ['Lexus RX 350', 'ian3@gt.edu'],
         ['Tesla Model X', 'jb@123.com'],
         ['Ford Pinto', 'michael@gt.edu'],
         ['Pontiac Firebird', 'jb@123.com'],
         ['VW Golf', 'michael@gt.edu'],
         ['Kia Sportage', 'michael@gt.edu'],
         ['Toyota RAV 4', 'jb@123.com']
    ]
);




/*
Pushpins
input = [[url, description, corkboard_title, user_id], ...]
Sadly, we can't randomize much here if we want the data to make sense
*/
CREATE OR REPLACE FUNCTION create_pushpins(text[][]) RETURNS void AS $$
DECLARE
  _pushpin text[];
BEGIN
  FOREACH _pushpin SLICE 1 IN ARRAY $1
        LOOP
            INSERT INTO pushpin (
                url,
                description,
                corkboard_id,
                        user_id,
                        created_at
            )
                VALUES (
                    _pushpin[1],
                    _pushpin[2],
                        (SELECT id FROM (
                            SELECT title, id
                            FROM private_corkboard
                            UNION
                            SELECT title, id
                            FROM public_corkboard
                            ) AS all_corkboards WHERE all_corkboards.title = _pushpin[3]),
                        (SELECT id FROM "user" WHERE email = _pushpin[4]),
                        (SELECT timestamp '2018-01-01' + random() * (timestamp '2018-01-31' - timestamp '2018-01-01'))
                );
  END LOOP;
END;
$$ LANGUAGE plpgsql;


SELECT 'Creating Public Corkboard Pushpins...' FROM create_pushpins(
    ARRAY[
                        ['https://upload.wikimedia.org/wikipedia/en/1/1b/Tranquility_Base_Hotel_%26_Casino_%28AM%29.jpg',
                         'One of my faves!',
                         'Tranquility Base Hotel & Casino',
                         'ian3@gt.edu'],
                        ['https://upload.wikimedia.org/wikipedia/en/9/91/The_1975_-_I_Like_It_When_You_Sleep%2C_for_You_Are_So_Beautiful_Yet_So_Unaware_of_It.png',
                         'These guys are _amazing_ live!',
                         'I Like It When You Sleep',
                         'michael@gt.edu'],
                        ['https://pictures.dealer.com/d/delrayhondavtg/0103/386effab492fae916e569403ef401c50x.jpg',
                         'Hawt Red Car!',
                         'BMW 328i',
                         'jb@123.com']
     ]
);


/*
Comments
input = [an, array, of, comment_text, strings]
Comment creator, date, and corkboard will be picked at random
@TODO: Limit the comments to users allowed to comment so our seed data is logical
*/
CREATE OR REPLACE FUNCTION create_comments(text[][]) RETURNS void AS $$
DECLARE
  _comment text;
BEGIN
  FOREACH _comment IN ARRAY $1
	LOOP
		INSERT INTO comment (
            comment_text,
            created_at,
            user_id,
            pushpin_id
		)
		VALUES (
		    _comment,
		    (SELECT timestamp '2017-10-01' + random() * (timestamp '2017-10-31' - timestamp '2017-10-01')),
			(SELECT id FROM "user" ORDER BY RANDOM() LIMIT 1 ),
			(SELECT id FROM pushpin ORDER BY RANDOM() LIMIT 1)

		);
  END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT 'Creating Comments...' FROM create_comments(
    ARRAY[
        'What makes you think she is a witch?',
        'A scratch? Your arm’s off!',
        'You must cut down the mightiest tree in the forest with…a herring!',
        'I’ve soiled my armour!',
        'What... is the air-speed velocity of an unladen swallow?',
        'On second thought, lets not go to Camelot. Tis a silly place.',
        'Strange women lying in ponds distributing swords is no basis for a system of government!',
        'I fart in your general direction. Your mother was a hamster and your father smelt of elderberries!',
        'I am an enchanter. There are some who call me... Tim?!'
    ]
);

/*
Likes
input = [[user_email, pin_title]]
@TODO: check validity of likes (can like pushpins on boards they can see)
Note-- doing lookups on the pin title is NOT WISE!
For seed data we know that to be unique so this works, but don't do it that way with user input
*/
CREATE OR REPLACE FUNCTION create_likes(text[]) RETURNS void AS $$
DECLARE
  _like text[];
BEGIN
FOREACH _like SLICE 1 IN ARRAY $1
LOOP
  INSERT INTO "like" (
    user_id,
    pushpin_id
  ) VALUES (
    (SELECT id FROM "user" WHERE email = _like[1]),
    (SELECT id FROM pushpin ORDER BY RANDOM() LIMIT 1)
  );
  END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT 'Creating Likes...' FROM create_likes(
    ARRAY[
        ['ian3@gt.edu'],
        ['michael@gt.edu'],
        ['jb@123.com']
    ]
);


/*
Tags
input = [array,of,tags]
pushpins selected randomly
*/
CREATE OR REPLACE FUNCTION create_tags(text[]) RETURNS void AS $$
DECLARE
  _tag text;
BEGIN
  FOREACH _tag IN ARRAY $1
	LOOP
	    INSERT INTO "tag" (
	        text,
	        pushpin_id
	    )
		VALUES (
		    _tag,
		    (SELECT id FROM pushpin ORDER BY RANDOM() LIMIT 1)
		);
  END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT 'Creating Tags...' FROM create_tags(
    ARRAY[
        'Happy',
        'To research',
        'To buy',
        'Things that scare me',
        'Favorites'
    ]
);

/*
Watches
*/
CREATE OR REPLACE FUNCTION create_watches(text[][]) RETURNS void AS $$
DECLARE
  _watch text[];
BEGIN
  FOREACH _watch SLICE 1 IN ARRAY $1
	LOOP
	    INSERT INTO watch (
	        user_id,
	        corkboard_id
	    )
		VALUES (
		    (SELECT id FROM "user" WHERE email = _watch[1]),
		    (SELECT id FROM (
			    SELECT title, id
			    FROM private_corkboard
			    UNION
			    SELECT title, id
			    FROM public_corkboard
			    ) AS all_corkboards WHERE all_corkboards.title = _watch[2])
		);
  END LOOP;
END;
$$ LANGUAGE plpgsql;

SELECT 'Creating Watches...' FROM create_watches(
    ARRAY[
        ['ian3@gt.edu', 'Icky Thump'],
        ['michael@gt.edu', 'Tesla Model X'],
        ['jb@123.com','Ford Pinto']
    ]
);
