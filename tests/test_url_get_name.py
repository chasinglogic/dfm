"""Test automatic name detection for URLs."""


from dfm.profile import get_name


def test_get_name():
    test_cases = [
        ("https://github.com/chasinglogic/dotfiles", "chasinglogic"),
        ("http://github.com/chasinglogic/dotfiles", "chasinglogic"),
        ("git@github.com:chasinglogic/dotfiles", "chasinglogic"),
        ("keybase://private/chasinglogic/dotfiles", "chasinglogic"),
    ]

    for url, expected in test_cases:
        assert get_name(url) == expected
