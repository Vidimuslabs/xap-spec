# xap-client (Python)

```bash
pip install -e .
```

```python
from xap_client import Client

c = Client("http://localhost:8080/xap/v1")
print(c.get_anchors())

admin = Client("http://localhost:8080/xap/v1", token="…")
admin.verify_chain()
```

stdlib only (`urllib`). See the parent [README](../README.md).
