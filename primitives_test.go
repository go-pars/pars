package pars_test

import (
	"strings"
	"testing"

	"github.com/ktnyt/pars"
	"github.com/stretchr/testify/require"
)

func TestEpsilon(t *testing.T) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	err := pars.Epsilon(s, pars.VoidResult)
	require.NoError(t, err)
}

func BenchmarkEpsilon(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Epsilon
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func BenchmarkDryEpsilon(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.Epsilon)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestHead(t *testing.T) {
	t.Run("matches empty string", func(t *testing.T) {
		s := pars.NewState(strings.NewReader(""))
		err := pars.Head(s, pars.VoidResult)
		require.NoError(t, err)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		s.Advance(1)
		err := pars.Head(s, pars.VoidResult)
		require.Error(t, err)
	})
}

func BenchmarkHead(b *testing.B) {
	s := pars.NewState(strings.NewReader(""))
	p := pars.Dry(pars.Head)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestEOF(t *testing.T) {
	t.Run("matches empty string", func(t *testing.T) {
		s := pars.NewState(strings.NewReader(""))
		err := pars.EOF(s, pars.VoidResult)
		require.NoError(t, err)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		err := pars.EOF(s, pars.VoidResult)
		require.Error(t, err)
	})
}

func BenchmarkEOF(b *testing.B) {
	s := pars.NewState(strings.NewReader(""))
	p := pars.Dry(pars.EOF)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestCut(t *testing.T) {
	s := pars.NewState(strings.NewReader("wikiwikiwiki"))
	e := "wiki"
	p := pars.String(e)
	r := pars.Result{}
	var err error
	for err = nil; err != nil; err = p(s, &r) {
		err := p(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)

		pars.Cut(s, &r)
		require.Equal(t, 0, s.Index)
		require.NotEqual(t, 0, s.Position)

		r = pars.Result{}
	}
}

func TestByte(t *testing.T) {
	t.Run("matches byte", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.Byte(e)(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := byte('h')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.Byte(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkByte(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.Byte('H'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotByte(t *testing.T) {
	t.Run("does not match byte", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotByte('h')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotByte(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotByte(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.NotByte('h'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestBytes(t *testing.T) {
	t.Run("matches bytes", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.Bytes('h', 'H')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.Bytes('h', 'g')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkBytes(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.Bytes('h', 'H'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotBytes(t *testing.T) {
	t.Run("matches bytes", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotBytes('g', 'G')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotBytes('h', 'H')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotBytes(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.NotBytes('g', 'G'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestByteRange(t *testing.T) {
	t.Run("matches byte range", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.ByteRange('A', 'Z')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.ByteRange('a', 'z')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkByteRange(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.ByteRange('A', 'Z'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotByteRange(t *testing.T) {
	t.Run("matches byte range", func(t *testing.T) {
		e := byte('H')
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotByteRange('a', 'z')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.NotByteRange('A', 'Z')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotByteRange(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.NotByteRange('a', 'z'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestByteSlice(t *testing.T) {
	t.Run("matches byte slice", func(t *testing.T) {
		e := []byte("Hello")
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.ByteSlice(e)(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := []byte("hello")
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.ByteSlice(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkByteSlice(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.ByteSlice([]byte("Hello")))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestRune(t *testing.T) {
	t.Run("matches rune", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.Rune(e)(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍖🍣"))
		r := pars.Result{}
		err := pars.Rune(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkRune(b *testing.B) {
	s := pars.NewState(strings.NewReader("🍺🍣"))
	p := pars.Dry(pars.Rune('🍺'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotRune(t *testing.T) {
	t.Run("matches rune", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.NotRune('🍖')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.NotRune(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotRune(b *testing.B) {
	s := pars.NewState(strings.NewReader("🍺🍣"))
	p := pars.Dry(pars.Rune('🍖'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestRunes(t *testing.T) {
	t.Run("matches rune", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.Runes('🍺', '🍣')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.Runes('🍖', '🍣')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkRunes(b *testing.B) {
	s := pars.NewState(strings.NewReader("🍺🍣"))
	p := pars.Dry(pars.Runes('🍺', '🍣'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotRunes(t *testing.T) {
	t.Run("matches rune", func(t *testing.T) {
		e := '🍺'
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.NotRunes('🍖', '🍣')(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("🍺🍣"))
		r := pars.Result{}
		err := pars.NotRunes('🍺', '🍣')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotRunes(b *testing.B) {
	s := pars.NewState(strings.NewReader("🍺🍣"))
	p := pars.Dry(pars.NotRunes('🍖', '🍣'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestRuneRange(t *testing.T) {
	t.Run("matches rune range", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("こんにちは"))
		r := pars.Result{}
		err := pars.RuneRange('あ', 'ん')(s, &r)
		require.NoError(t, err)
		require.IsType(t, 'こ', r.Value)
		require.Equal(t, 'こ', r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("こんにちは"))
		r := pars.Result{}
		err := pars.RuneRange('ア', 'ン')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkRuneRange(b *testing.B) {
	s := pars.NewState(strings.NewReader("こんにちは"))
	p := pars.Dry(pars.RuneRange('あ', 'ん'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestNotRuneRange(t *testing.T) {
	t.Run("matches rune range", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("こんにちは"))
		r := pars.Result{}
		err := pars.NotRuneRange('ア', 'ン')(s, &r)
		require.NoError(t, err)
		require.IsType(t, 'こ', r.Value)
		require.Equal(t, 'こ', r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		s := pars.NewState(strings.NewReader("こんにちは"))
		r := pars.Result{}
		err := pars.NotRuneRange('あ', 'ん')(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkNotRuneRange(b *testing.B) {
	s := pars.NewState(strings.NewReader("こんにちは"))
	p := pars.Dry(pars.NotRuneRange('ア', 'ン'))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestRuneSlice(t *testing.T) {
	t.Run("matches byte slice", func(t *testing.T) {
		e := []rune("こんにちは")
		s := pars.NewState(strings.NewReader("こんにちは🍣🍺"))
		r := pars.Result{}
		err := pars.RuneSlice(e)(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := []rune("🍣🍺")
		s := pars.NewState(strings.NewReader("こんにちは🍣🍺"))
		r := pars.Result{}
		err := pars.RuneSlice(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkRuneSlice(b *testing.B) {
	s := pars.NewState(strings.NewReader("こんにちは🍣🍺"))
	p := pars.Dry(pars.RuneSlice([]rune("こんにちは")))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}

func TestString(t *testing.T) {
	t.Run("matches string", func(t *testing.T) {
		e := "Hello"
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.String(e)(s, &r)
		require.NoError(t, err)
		require.IsType(t, e, r.Value)
		require.Equal(t, e, r.Value)
	})

	t.Run("returns error", func(t *testing.T) {
		e := "hello"
		s := pars.NewState(strings.NewReader("Hello world!"))
		r := pars.Result{}
		err := pars.String(e)(s, &r)
		require.Error(t, err)
		require.Nil(t, r.Value)
	})
}

func BenchmarkString(b *testing.B) {
	s := pars.NewState(strings.NewReader("Hello world!"))
	p := pars.Dry(pars.String("Hello"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p(s, pars.VoidResult)
	}
}
