CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fb_external_post
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();