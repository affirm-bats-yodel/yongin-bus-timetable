-- List of Buses
CREATE TABLE IF NOT EXISTS bus_lists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    create_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    name VARCHAR UNIQUE NOT NULL
);

-- a Timetable of buses
CREATE TABLE IF NOT EXISTS bus_timetables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bus_id UUID REFERENCES bus_lists(id), -- reference from bus_lists.id
    stop_name VARCHAR NOT NULL, -- bus stop name
    depart_at TIME NOT NULL, -- bus depart time
);
