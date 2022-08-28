-- CREATE schema words;
-- CREATE DATABASE words_permutation;
USE words_permutation;

-- DROP TABLE words;
-- DROP TABLE IF EXISTS words;

CREATE TABLE IF NOT EXISTS words (
    word VARCHAR(255) NOT NULL,
	permutation_table_key VARCHAR(255),
    PRIMARY KEY (word)
);

-- INSERT INTO words (word) VALUES  ("a"), ("aa"), ("aaa"), ("aah"), ("aahed"), ("aahing"), ("aahs"), ("aal"), ("aalii"), ("aaliis"), ("aals"), ("aam"), ("aardvark"), ("aardvarks"), ("aardwolf"), ("aardwolves"), ("aargh"), ("aaron"), ("aaronic"), ("aarrgh"), ("aarrghh"), ("aas"), ("aasvogel"), ("aasvogels"), ("ab"), ("aba"), ("abac"), ("abaca"), ("abacas"), ("abacate"), ("abacaxi"), ("abacay"), ("abaci"), ("abacinate"), ("abacination"), ("abacisci"), ("abaciscus"), ("abacist"), ("aback"), ("abacli"), ("abacot"), ("abacterial"), ("abactinal"), ("abactinally"), ("abaction"), ("abactor"), ("abaculi"), ("abaculus"), ("abacus"), ("abacuses");
-- UPDATE words
-- SET associated_words_table = 2
-- WHERE word="test";
-- INSERT INTO words (word) VALUES  ("a"), ("aa"), ("aaa"), ("aah"), ("aahed"), ("aahing"), ("aahs"), ("aal"), ("aalii"), ("aaliis");
SELECT * FROM `int` WHERE word != 'int';
-- INSERT INTO 'int' (word) VALUES ("int"), ("nit")
