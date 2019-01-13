
CREATE TABLE IF NOT EXISTS "user"(
  id SERIAL PRIMARY KEY,
  email VARCHAR(50) UNIQUE NOT NULL,
  pin INT NOT NULL,
  user_name VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS follow (
  user_id SERIAL REFERENCES "user"(id),
  followee_id SERIAL REFERENCES "user"(id),
  PRIMARY KEY (user_id, followee_id)
);

CREATE TABLE IF NOT EXISTS category (
  category_name VARCHAR(50) NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS corkboard (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS public_corkboard (
  id SERIAL REFERENCES corkboard(id),
  title VARCHAR(50) NOT NULL,
  user_id SERIAL REFERENCES "user"(id) NOT NULL,
  category_name VARCHAR(250) REFERENCES category(category_name),
  PRIMARY KEY (title, user_id)
);

CREATE TABLE IF NOT EXISTS private_corkboard (
  id SERIAL REFERENCES corkboard(id),
  title VARCHAR(50) NOT NULL,
  user_id SERIAL REFERENCES "user"(id) NOT NULL,
  category_name VARCHAR(250) REFERENCES category(category_name),
  password VARCHAR(50) NOT NULL,
  PRIMARY KEY (title, user_id)
);

CREATE TABLE IF NOT EXISTS pushpin (
  id SERIAL PRIMARY KEY,
  url VARCHAR(1024),
  description VARCHAR(255),
  created_at TIMESTAMP,
  corkboard_id SERIAL REFERENCES corkboard(id),
  user_id SERIAL REFERENCES "user"(id)
);

CREATE TABLE IF NOT EXISTS comment (
  id SERIAL PRIMARY KEY,
  comment_text VARCHAR(250) NOT NULL,
  created_at TIMESTAMP,
  user_id SERIAL REFERENCES "user"(id),
  pushpin_id SERIAL REFERENCES pushpin(id)
);

CREATE TABLE IF NOT EXISTS watch (
  user_id SERIAL REFERENCES "user"(id),
  corkboard_id SERIAL REFERENCES corkboard(id) ,
  PRIMARY KEY (user_id, corkboard_id)
);


CREATE TABLE IF NOT EXISTS "like" (
  user_id SERIAL REFERENCES "user"(id),
  pushpin_id SERIAL REFERENCES pushpin(id),
  PRIMARY KEY (user_id, pushpin_id)
);

CREATE TABLE IF NOT EXISTS tag (
  text VARCHAR(20) NOT NULL,
  pushpin_id SERIAL REFERENCES pushpin(id),
  PRIMARY KEY (text, pushpin_id)
);


SELECT public_corkboard.category_name, public_corkboard.title, COUNT(watch.*), pushpin.url, MAX(pushpin.created_at)
FROM public_corkboard
INNER JOIN pushpin on pushpin.id = public_corkboard.id
INNER JOIN "user" on public_corkboard.user_id = "user".id
LEFT JOIN watch
	ON (
		watch.corkboard_id = public_corkboard.id
  )
where public_corkboard.id = 1
GROUP BY public_corkboard.category_name, public_corkboard.title, pushpin.url;


SELECT 
	pushpin.description, 
	public_corkboard.title,  
	"user".user_name
FROM "user"
INNER JOIN public_corkboard 
ON ("user".id = public_corkboard.user_id)
INNER JOIN pushpin 
ON (pushpin.corkboard_id = public_corkboard.id)
JOIN tag 
ON (tag.pushpin_id = pushpin.id)
JOIN category 
ON (category.category_name = public_corkboard.category_name)
WHERE pushpin.description LIKE '%input_search_term%'
OR tag.text LIKE '%input_search_term%'
OR category.category_name LIKE '%input_search_term%'
GROUP BY pushpin.description, public_corkboard.title, "user".user_name;
