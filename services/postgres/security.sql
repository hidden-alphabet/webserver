DO $$
BEGIN

REVOKE ALL ON user.meta FROM public;
REVOKE ALL ON user.account FROM public;
REVOKE ALL ON user.session FROM public;

END;
$$;
