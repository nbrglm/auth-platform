# flake8: noqa

if __import__("typing").TYPE_CHECKING:
    # import apis into api package
    from auth_platform_py_sdk.api.auth_api import AuthApi
    
else:
    from lazy_imports import LazyModule, as_package, load

    load(
        LazyModule(
            *as_package(__file__),
            """# import apis into api package
from auth_platform_py_sdk.api.auth_api import AuthApi

""",
            name=__name__,
            doc=__doc__,
        )
    )
