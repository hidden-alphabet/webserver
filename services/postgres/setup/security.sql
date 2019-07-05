DO $$
BEGIN

REVOKE ALL ON web.api FROM public;
REVOKE ALL ON web.contact FROM public;
REVOKE ALL ON web.account FROM public;
REVOKE ALL ON web.session FROM public;

END;
$$;
