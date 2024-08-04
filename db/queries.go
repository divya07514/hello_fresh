package db

var CreateDb = "CREATE DATABASE IF NOT EXISTS hello_fresh"
var InsertRecipe = "INSERT INTO recipe_list (postcode, recipe, `from`, `to`, day_of_week) VALUES (?, ?, ?, ?, ?)"
var UniqueRecipeCount = "select count(distinct recipe) from recipe_list"
var UniqueRecipeAndCount = "SELECT COUNT(*) AS recipe_count, recipe FROM recipe_list GROUP BY recipe ORDER BY recipe ASC"
var BusiestPostCode = "select count(postcode), postcode from recipe_list group by postcode order by count(postcode) DESC limit 1"
var DeliveriesToPostcode = "select count(*) from recipe_list where postcode =? and `from` >=? and `to` <=?"
var RecipesLike = "SELECT DISTINCT(recipe) FROM recipe_list WHERE %s ORDER BY recipe ASC"
var CreateTable = "CREATE TABLE IF NOT EXISTS hello_fresh.recipe_list (id INT UNSIGNED NOT NULL AUTO_INCREMENT, postcode VARCHAR(255) DEFAULT NULL,recipe VARCHAR(255) DEFAULT NULL, `from` INT DEFAULT NULL, `to` INT DEFAULT NULL, day_of_week VARCHAR(255) DEFAULT NULL, PRIMARY KEY (id), KEY `postcode_idx` (`postcode`), KEY `recipe_idx` (`recipe`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
