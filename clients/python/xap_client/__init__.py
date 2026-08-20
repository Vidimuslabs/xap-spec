"""XAP HTTP client — thin wrapper over a self-hosted node (/xap/v1)."""

from .client import Client, XAPError

__all__ = ["Client", "XAPError"]
__version__ = "0.1.0"
