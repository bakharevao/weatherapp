CREATE TABLE IF NOT EXISTS city_coordinates (
    city      VARCHAR(100) NOT NULL PRIMARY KEY,
    latitude  DOUBLE       NOT NULL,
    longitude DOUBLE       NOT NULL
    );

INSERT INTO city_coordinates (city, latitude, longitude) VALUES
    ('limassol', 34.7071, 33.0226),
    ('nicosia', 35.1856, 33.3823),
    ('moscow', 55.7558, 37.6173),
    ('paris', 48.8566, 2.3522),
    ('dubai', 25.2048, 55.2708)
ON DUPLICATE KEY UPDATE
    latitude = VALUES(latitude),
    longitude = VALUES(longitude);