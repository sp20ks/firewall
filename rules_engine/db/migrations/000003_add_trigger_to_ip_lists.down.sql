DROP TRIGGER IF EXISTS check_ip_list_conflict ON resource_ip_list;
DROP FUNCTION IF EXISTS prevent_conflicting_ip_lists;
DROP INDEX IF EXISTS unique_ip_per_resource;
