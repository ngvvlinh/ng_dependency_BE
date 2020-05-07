CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fb_external_conversation
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fb_external_comment
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();

CREATE TRIGGER notify_pgrid AFTER INSERT ON history.fb_external_message
  FOR EACH ROW EXECUTE PROCEDURE notify_pgrid_id();
