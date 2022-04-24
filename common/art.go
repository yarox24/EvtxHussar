package common

import (
	"fmt"
	"strings"
)

func Hussar_art(VERSION string, URL string) {
	asciiart := "\n         :.                             \n" +
		"         $$$$*                          \n" +
		"         V$$$$$   N$M                   \n" +
		"         :$$$$$* *$$$$   :              \n" +
		"          $$$$$: $$$$$  V$$             \n" +
		"          $$$$$ V$$$$  I$$$:            \n" +
		"          $$$$* $$$N  N$$$N  V$M        \n" +
		"          M$$N N$$I  $$$$: V$$$$        \n" +
		"          *$$ :$$V *$$$* *$$$$$         \n" +
		"           :. M$* V$$V .$$$$M:  V$$:    \n" +
		"          VIF*.. *$V  N$$M:  V$$$$$     \n" +
		"               :IV  M$M:..F$$$$$$I      \n" +
		"                  *F  :M$$$I*.    :VN   \n" +
		"                    N.:.   .*I$$$$$$*   \n" +
		"  EvtxHussar VERSION $ *$$$NIV*:.       \n" +
		"                     I*      .::*VVM$*  \n" +
		"  URL $$NNN*    \n" +
		"                     FV .............   \n" +
		"                     FV $$$$$$$$$$$V    \n" +
		"    Provided by      FV                 \n" +
		"  Jarosław Oparka    FV.$$$$$$$$$$N:    \n" +
		"                     FV                 \n" +
		"                     FV.$$$$$$$$$$V     \n" +
		"                     FV                 \n" +
		"             .*VFVV* FV.$$$$$$$$$$.     \n" +
		"         .FF:        FV ::::::::        \n" +
		"       VM.           FV.MNMMNMMNMV      \n" +
		"     VN              FV *******:        \n" +
		"    N*               FV IIIIIMMMM  *    \n" +
		"   $:                               M   \n" +
		"  NV           :I$$$$$NV.:MNNNNN.   .V  \n" +
		" :$          *$$$$$$$$$$$M           $  \n" +
		" IM         *$$$$$$$$$$$$$$ $$$$*    $. \n" +
		" NF         $$$$$$$$$$$$$$$V         N: \n" +
		" V$         N$$$$$$$$$$$$$$$*$$V     $  \n" +
		"  $                 F  :$$$M         $  \n" +
		"  VN               N$I  V$$M        $.  \n" +
		"   FI              $$M   :$$       M:   \n" +
		"    *N              I            .N     \n" +
		"      IF                        M*      \n" +
		"        VN:                  *N*        \n" +
		"           *MV:         .*FF:           \n" +
		"                .::**:.                 "

	// Insert Version
	asciiart = strings.Replace(asciiart, "VERSION", fmt.Sprintf("%-7s", VERSION), 1)
	asciiart = strings.Replace(asciiart, "URL", fmt.Sprintf("%s", URL), 1)

	fmt.Println(asciiart)
}
