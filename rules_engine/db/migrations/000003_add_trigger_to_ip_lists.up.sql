CREATE UNIQUE INDEX unique_ip_per_resource ON resource_ip_list (resource_id, ip_list_id);
CREATE OR REPLACE FUNCTION prevent_conflicting_ip_lists()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM resource_ip_list ril
        JOIN ip_lists il1 ON ril.ip_list_id = il1.id
        JOIN ip_lists il2 ON il1.ip = il2.ip
        WHERE ril.resource_id = NEW.resource_id
        AND il1.list_type <> il2.list_type
    ) THEN
        RAISE EXCEPTION 'IP already exists in the opposite list type for this resource';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_ip_list_conflict
BEFORE INSERT OR UPDATE ON resource_ip_list
FOR EACH ROW EXECUTE FUNCTION prevent_conflicting_ip_lists();
