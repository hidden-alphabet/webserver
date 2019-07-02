DO $$
BEGIN

REVOKE ALL ON web.api FROM public;
REVOKE ALL ON web.user FROM public;
REVOKE ALL ON web.resource FROM public;

END;
$$;
