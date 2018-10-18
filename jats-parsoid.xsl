<?xml version="1.0"?>

<xsl:stylesheet version="1.0"
    xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
    xmlns:xlink="http://www.w3.org/1999/xlink"
    xmlns:mml="http://www.w3.org/1998/Math/MathML"
    exclude-result-prefixes="xlink mml">

    <xsl:strip-space elements="*"/>

    <xsl:variable name="vLower" select="'abcdefghijklmnopqrstuvwxyz'"/>
    <xsl:variable name="vUpper" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ'"/>

    <xsl:output method="html" doctype-system="about:legacy-compat"/>

    <xsl:template match="/">
        <html>
            <xsl:call-template name="make-html-header"/>
            <body>
                <xsl:apply-templates/>
            </body>
        </html>
    </xsl:template>

    <xsl:template name="make-html-header">
        <head>
            <title>
                <xsl:text>Converted JATS paper:</xsl:text>
            </title>
        </head>
    </xsl:template>

    <xsl:template match="article">
        <xsl:call-template name="make-article"/>
    </xsl:template>


    <xsl:template name="make-article">
        <xsl:apply-templates select="front"/>
        <xsl:for-each select="body">
            <h2 id="Paper"><xsl:text>Paper</xsl:text></h2>
            <xsl:apply-templates/>
        </xsl:for-each>
        <xsl:apply-templates select="back"/>
    </xsl:template>

    <xsl:template match="front">
        <section>
            <h2 id="Journal_Information"><xsl:text>Journal Information</xsl:text></h2>
            <xsl:for-each select="journal-meta">
                <xsl:apply-templates/>
            </xsl:for-each> <!-- journal-meta -->
        </section>

        <section>
            <xsl:for-each select="article-meta">
                <xsl:apply-templates/>
            </xsl:for-each> <!-- article-meta -->
        </section>
    </xsl:template>

    <xsl:template match="back">
        <xsl:apply-templates select="ack | ref-list"/>
    </xsl:template>


    <!-- journal-meta -->

    <xsl:template match="journal-meta/*" priority="0.0">
        <!-- pass -->
    </xsl:template>

    <xsl:template match="journal-title">
        <p>
            <xsl:text>Title: </xsl:text>
            <xsl:value-of select="."/>
        </p>
    </xsl:template>

    <xsl:template match="journal-title-group">
        <xsl:apply-templates/>
    </xsl:template>

    <!-- article-meta -->

    <xsl:template match="article-meta/*" priority="0.0">
        <!-- pass -->
    </xsl:template>

    <xsl:template match="title-group/article-title">
        <xsl:param name="contents">
            <xsl:apply-templates/>
        </xsl:param>
        <h1>
            <xsl:attribute name="id">
                <xsl:value-of select="translate($contents, ' ', '_')"/>
            </xsl:attribute>
            <xsl:copy-of select="$contents"/>
        </h1>
    </xsl:template>

    <xsl:template match="alt-title">
        <xsl:param name="contents">
            <xsl:apply-templates/>
        </xsl:param>
        <p>
            <xsl:text>Alternative Title: </xsl:text>
            <xsl:copy-of select="$contents"/>
        </p>
    </xsl:template>

    <xsl:template match="title-group">
        <xsl:apply-templates/>
    </xsl:template>

    <xsl:template match="role">
        <xsl:text> (</xsl:text>
            <xsl:value-of select="."/>
        <xsl:text>)</xsl:text>
    </xsl:template>

    <xsl:template name="contrib-identify">
        <li>
            <xsl:for-each select="anonymous | collab | collab-alternatives/* | name | name-alternatives/*">
                <xsl:apply-templates select="."/>
            </xsl:for-each>
            <xsl:apply-templates select="role"/>
        </li>
    </xsl:template>

    <xsl:template name="contrib-info">
        <xsl:apply-templates mode="metadata" select="aff"/>
    </xsl:template>

    <!-- <xsl:template match="role">
        <p><xsl:value-of select="."/></p>
    </xsl:template> -->

    <xsl:template match="article-meta/contrib-group">
        <ul>
            <xsl:for-each select="contrib">
                <xsl:call-template name="contrib-identify"/>
                <xsl:call-template name="contrib-info"/>
            </xsl:for-each> <!-- contrib -->
        </ul>
    </xsl:template>

    <xsl:template match="aff">
        <p>
            <xsl:choose>
                <xsl:when test="addr-line">
                    <xsl:if test="label">
                        <xsl:text> [</xsl:text>
                            <xsl:value-of select="label"/>
                        <xsl:text>]</xsl:text>
                    </xsl:if>
                    <xsl:value-of select="addr-line"/>
                </xsl:when>
                <xsl:otherwise>
                    <xsl:value-of select="."/>
                </xsl:otherwise>
            </xsl:choose>
        </p>
    </xsl:template>

    <xsl:template match="pub-date">
        <p>
            <xsl:text>Publication date</xsl:text>
            <xsl:if test="@pub-type">
                <xsl:text> (</xsl:text>
                    <xsl:value-of select="@pub-type"/>
                <xsl:text>)</xsl:text>
            </xsl:if>
            <xsl:text>: </xsl:text>
            <xsl:value-of select="month"/>
            <xsl:text>/</xsl:text>
            <xsl:value-of select="year"/>
        </p>
    </xsl:template>

    <xsl:template match="abstract">
        <xsl:param name="abstype">
            <xsl:value-of select="@abstract-type"/>
        </xsl:param>
        <xsl:param name="caps">
            <xsl:value-of select="concat(translate(substring($abstype,1,1), $vLower, $vUpper), substring($abstype, 2))"/>
        </xsl:param>
        <section>
            <xsl:choose>
                <xsl:when test="@abstract-type">
                    <h2>
                        <xsl:attribute name="id">
                            <xsl:value-of select="$caps"/>
                        </xsl:attribute>
                        <xsl:value-of select="$caps"/>
                    </h2>
                </xsl:when>
                <xsl:otherwise>
                    <h2 id="Abstract">Abstract</h2>
                </xsl:otherwise>
            </xsl:choose>
            <xsl:apply-templates/>
        </section>
    </xsl:template>

    <!-- body -->

    <!-- back -->

    <xsl:template match="ack">
        <section>
            <h2 id="Acknowledgements"><xsl:text>Acknowledgements</xsl:text></h2>
            <xsl:apply-templates/>
        </section>
    </xsl:template>

    <xsl:template match="ref-list">
        <section>
            <h2 id="References"><xsl:text>References</xsl:text></h2>
            <ul>
                <xsl:for-each select="ref">
                    <xsl:call-template name="ref-info"/>
                </xsl:for-each> <!-- ref -->
            </ul>
        </section>
    </xsl:template>

    <xsl:template name="ref-info">
        <li>
            <xsl:if test="label">
                <xsl:text>[</xsl:text>
                    <xsl:value-of select="label"/>
                <xsl:text>] </xsl:text>
            </xsl:if>
            <xsl:apply-templates select="element-citation"/>
        </li>
    </xsl:template>


    <!-- some basic formatting -->

    <xsl:template match="bold">
        <b><xsl:apply-templates/></b>
    </xsl:template>

    <xsl:template match="italic">
        <i><xsl:apply-templates/></i>
    </xsl:template>

    <xsl:template match="monospace">
        <tt><xsl:apply-templates/></tt>
    </xsl:template>

  <!-- ============================================================= -->
  <!--  REGULAR (DEFAULT) MODE                                       -->
  <!-- ============================================================= -->


  <xsl:template match="sec">
    <section>
      <!-- <xsl:call-template name="named-anchor"/>-->
      <xsl:apply-templates select="title"/>
      <!-- <xsl:apply-templates select="sec-meta"/> -->
      <xsl:apply-templates mode="drop-title"/>
    </section>
  </xsl:template>

  <xsl:template match="*" mode="drop-title">
    <xsl:apply-templates select="."/>
  </xsl:template>

  <xsl:template match="title | sec-meta" mode="drop-title"/>

  <xsl:template match="p | license-p">
    <p>
      <xsl:apply-templates/>
    </p>
  </xsl:template>

    <!-- ============================================================= -->
    <!--  Writing a name                                               -->
    <!-- ============================================================= -->

    <!-- Called when displaying structured names in metadata         -->

  <xsl:template match="name">
    <xsl:apply-templates select="prefix" mode="inline-name"/>
    <xsl:apply-templates select="surname[../@name-style='eastern']"
      mode="inline-name"/>
    <xsl:apply-templates select="given-names" mode="inline-name"/>
    <xsl:apply-templates select="surname[not(../@name-style='eastern')]"
      mode="inline-name"/>
    <xsl:apply-templates select="suffix" mode="inline-name"/>
  </xsl:template>


  <xsl:template match="prefix" mode="inline-name">
    <xsl:apply-templates/>
    <xsl:if test="../surname | ../given-names | ../suffix">
      <xsl:text> </xsl:text>
    </xsl:if>
  </xsl:template>


  <xsl:template match="given-names" mode="inline-name">
    <xsl:apply-templates/>
    <xsl:if test="../surname[not(../@name-style='eastern')] | ../suffix">
      <xsl:text> </xsl:text>
    </xsl:if>
  </xsl:template>


  <xsl:template match="contrib/name/surname" mode="inline-name">
    <xsl:apply-templates/>
    <xsl:if test="../given-names[../@name-style='eastern'] | ../suffix">
      <xsl:text> </xsl:text>
    </xsl:if>
  </xsl:template>


  <xsl:template match="surname" mode="inline-name">
    <xsl:apply-templates/>
    <xsl:if test="../given-names[../@name-style='eastern'] | ../suffix">
      <xsl:text> </xsl:text>
    </xsl:if>
  </xsl:template>


  <xsl:template match="suffix" mode="inline-name">
    <xsl:apply-templates/>
  </xsl:template>


  <!-- string-name elements are written as is -->

  <xsl:template match="string-name">
    <xsl:apply-templates/>
  </xsl:template>


  <xsl:template match="string-name/*">
    <xsl:apply-templates/>
  </xsl:template>



<!-- ============================================================= -->
<!--  Figures, lists and block-level objectS                       -->
<!-- ============================================================= -->

  <xsl:template match="title">
    <xsl:if test="normalize-space(string(.))">
      <h3>
        <xsl:attribute name="id">
            <xsl:value-of select="translate(., ' ', '_')"/>
        </xsl:attribute>
        <xsl:apply-templates/>
      </h3>
    </xsl:if>
  </xsl:template>


  <xsl:template match="subtitle">
    <xsl:if test="normalize-space(string(.))">
      <h5>
        <xsl:attribute name="id">
            <xsl:value-of select="translate(., ' ', '_')"/>
        </xsl:attribute>
        <xsl:apply-templates/>
      </h5>
    </xsl:if>
  </xsl:template>

  <xsl:template match="aff/label | corresp/label | chem-struct/label |
    element-citation/label | mixed-citation/label | citation/label">
    <!-- these labels appear in line -->
    <span class="generated">[</span>
    <xsl:apply-templates/>
    <span class="generated">] </span>
  </xsl:template>



</xsl:stylesheet>
