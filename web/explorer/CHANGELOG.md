# CHANGELOG

## 1.8.0 (Dec 29, 2018)

- Add v3 blocks api (support pagination) ([1156c6f](https://gitlab.com/thunderdb/aletheia/commit/1156c6f))
- Block component support queries pagination ([9141282](https://gitlab.com/thunderdb/aletheia/commit/9141282))
- Bump version to 1.7.1 ([b217ab2](https://gitlab.com/thunderdb/aletheia/commit/b217ab2))
- Component BlockList use v3 api ([1508047](https://gitlab.com/thunderdb/aletheia/commit/1508047))
- Fix nginx config envsubst bug ([b71d0df](https://gitlab.com/thunderdb/aletheia/commit/b71d0df))
- Home component latest blocks use v3 api ([eddc71e](https://gitlab.com/thunderdb/aletheia/commit/eddc71e))
- Remove timezone in human readable time format ([f209718](https://gitlab.com/thunderdb/aletheia/commit/f209718))
- Upgrade node-sass ([1e7fae5](https://gitlab.com/thunderdb/aletheia/commit/1e7fae5))

## 1.7.1 (Nov 7, 2018)

- Fix nginx config envsubst bug.

## 1.7.0 (Nov 6, 2018)

- Dockerize.

## 1.6.1 (Sep 28, 2018)

- Fix block navigation render collision.

## 1.6.0 (Sep 28, 2018)

- Show errors explicitly in components;
- Add single page error card component;
- Fix incorrect URL jump while page loading;
- Catch errors while gathering SQL statements.

## 1.5.0 (Sep 27, 2018)

- Show error information in block lists.

## 1.4.0 (Sep 27, 2018)

- API upgrade, use continuous field count instead of height;
- Database selector component enhancement.

## 1.3.1 (Sep 3, 2018)

- Fix block list height start/end error.

## 1.3.0 (Sep 3, 2018)

- Add BlockList to list of blocks of the specified DB;
- Add MaxBlockHeight API, show latest block list at home page from
  newer to older;
- Add a SQL section at Block page to show relates SQL queries in this
  block;
- Colorize SQL types text.

## 1.2.0 (Aug 31, 2018)

- Add logo.

## 1.1.0 (Aug 27, 2018)

- Add a default block list to show latest blocks;
- Remember the last database address the user accessed;
- Realtime search on typing, no need to press the enter key now;
- Add ship script.

## 1.0.0 (Aug 26, 2018)

- Preview release.
