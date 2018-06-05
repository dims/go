// This file contains HMAC portability wrappers.

#include "goboringcrypto.h"

// Not in OpenSSL 1.1.  However, HMAC_CTX_copy expects an initialized
// target in OpenSSL 1.1.
int
_goboringcrypto_HMAC_CTX_copy_ex(GO_HMAC_CTX *dest, const GO_HMAC_CTX *src)
{
  // HMAC_CTX_copy lacks the const qualifier for the second parameter.
  return HMAC_CTX_copy(dest, (GO_HMAC_CTX *) src);
}
