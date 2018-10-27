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

    <xsl:template match="ext-link">
        <a rel="mw:ExtLink">
            <xsl:attribute name="href">
                <xsl:value-of select="@xlink:href"/>
            </xsl:attribute>
            <xsl:apply-templates/>
        </a>
    </xsl:template>

    <xsl:template match="xref[@ref-type='bibr']">
        <xsl:param name="rid">
            <xsl:value-of select="@rid"/>
        </xsl:param>
        <sup
            class="mw-ref"
            rel="dc:references"
            typeof="mw:Extension/ref"
        >
            <xsl:attribute name="id">
                <xsl:value-of select="concat('cite_ref-', @rid)"/>
            </xsl:attribute>
            <xsl:attribute name="data-mw">
                <xsl:text>{"name":"ref","attrs":{"name":"</xsl:text>
                <xsl:value-of select="@rid"/>
                <xsl:text>"},"body":{"id":"mw-reference-text-cite_note-</xsl:text>
                <xsl:value-of select="@rid"/>
                <xsl:text>"}}</xsl:text>
            </xsl:attribute>
            <a>
                <xsl:attribute name="href">
                    <xsl:value-of select="concat('#cite_note-', @rid)"/>
                </xsl:attribute>
                <span class="mw-reflink-text">
                    <xsl:text>[</xsl:text>
                    <xsl:apply-templates select="//article/back/ref-list/ref[@id=$rid]" mode="inline"/>
                    <xsl:text>]</xsl:text>
                </span>
            </a>
        </sup>
    </xsl:template>

    <xsl:template match="ref" mode="inline">
        <xsl:param name="offset">
            <xsl:number count="ref" from="ref-list" level="any"/>
        </xsl:param>
        <xsl:value-of select="$offset"/>
    </xsl:template>


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
            <div
                class="mw-references-wrap"
                typeof="mw:Extension/references"
            >
                <xsl:attribute name="data-mw">
                    <xsl:text>{"name":"references","attrs":{}}</xsl:text>
                </xsl:attribute>
                <ol class="mw-references references">
                    <xsl:for-each select="ref">
                        <xsl:call-template name="ref-info"/>
                    </xsl:for-each> <!-- ref -->
                </ol>
            </div>
        </section>
    </xsl:template>

    <xsl:template name="ref-info">
        <li>
            <xsl:attribute name="id">
                <xsl:value-of select="concat('cite_note-', @id)"/>
            </xsl:attribute>
            <a rel="mw:referenceBy">
                <xsl:attribute name="href">
                    <xsl:value-of select="concat('#cite_ref-', @id)"/>
                </xsl:attribute>
                <span class="mw-linkback-text"><xsl:text>↑ </xsl:text></span>
            </a>
            <span
                class="mw-reference-text"
            >
                <xsl:attribute name="id">
                    <xsl:value-of select="concat('mw-reference-text-cite_note-', @id)"/>
                </xsl:attribute>
                <cite
                    class="citation XXXXX"
                    typeof="mw:Transclusion"
                >
                    <xsl:attribute name="data-mw">
                        <xsl:apply-templates select="element-citation | mixed-citation" mode="data-mw"/>
                    </xsl:attribute>
                    <xsl:apply-templates select="element-citation | mixed-citation"/>
                </cite>
                <span>
                </span>
                <style
                    typeof="mw:Extension/templatestyles"
                >
                    <xsl:attribute name="data-mw">
                        <xsl:text>{"name":"templatestyles","attrs":{"src":"Module:Citation/CS1/styles.css"},"body":{"extsrc":""}}</xsl:text>
                    </xsl:attribute>
                    .mw-parser-output cite.citation{font-style:inherit}.mw-parser-output q{quotes:"\"""\"""'""'"}.mw-parser-output code.cs1-code{color:inherit;background:inherit;border:inherit;padding:inherit}.mw-parser-output .cs1-lock-free a{background:url("//upload.wikimedia.org/wikipedia/commons/thumb/6/65/Lock-green.svg/9px-Lock-green.svg.png")no-repeat;background-position:right .1em center}.mw-parser-output .cs1-lock-limited a,.mw-parser-output .cs1-lock-registration a{background:url("//upload.wikimedia.org/wikipedia/commons/thumb/d/d6/Lock-gray-alt-2.svg/9px-Lock-gray-alt-2.svg.png")no-repeat;background-position:right .1em center}.mw-parser-output .cs1-lock-subscription a{background:url("//upload.wikimedia.org/wikipedia/commons/thumb/a/aa/Lock-red-alt-2.svg/9px-Lock-red-alt-2.svg.png")no-repeat;background-position:right .1em center}.mw-parser-output .cs1-subscription,.mw-parser-output .cs1-registration{color:#555}.mw-parser-output .cs1-subscription span,.mw-parser-output .cs1-registration span{border-bottom:1px dotted;cursor:help}.mw-parser-output .cs1-hidden-error{display:none;font-size:100%}.mw-parser-output .cs1-visible-error{font-size:100%}.mw-parser-output .cs1-subscription,.mw-parser-output .cs1-registration,.mw-parser-output .cs1-format{font-size:95%}.mw-parser-output .cs1-kern-left,.mw-parser-output .cs1-kern-wl-left{padding-left:0.2em}.mw-parser-output .cs1-kern-right,.mw-parser-output .cs1-kern-wl-right{padding-right:0.2em}
                </style>
                <span>
                </span>
            </span>
        </li>
    </xsl:template>


    <!-- generate the mediawiki embedded data -->

    <xsl:template match="element-citation[@publication-type='other'] | mixed-citation[@publication-type='other']" mode="data-mw">
        <xsl:text>{"parts":[{"template":{"target":{"wt":"cite report","href":"./Template:Cite_report"},"params":{</xsl:text>
        <xsl:apply-templates select="article-title | year | collab | source | volume | publisher-name | fpage" mode="data-mw"/>
        <xsl:text>"noop":{"wt":"noop"}</xsl:text>
        <xsl:text>},"i":0}}]}</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation[@publication-type='journal'] | mixed-citation[@publication-type='journal']" mode="data-mw">
        <!-- {"parts":[{"template":{"target":{"wt":"cite journal","href":"./Template:Cite_journal"},"params":{"vauthors":{"wt":"foo, bar"},"date":{"wt":"2009"},"title":{"wt":"wibble"},"journal":{"wt":"blah"},"publisher":{"wt":"bastards"},"volume":{"wt":"12"}},"i":0}}]} -->
        <xsl:text>{"parts":[{"template":{"target":{"wt":"cite journal","href":"./Template:Cite_journal"},"params":{</xsl:text>
        <xsl:apply-templates select="article-title | year | collab | source | volume | publisher-name | fpage" mode="data-mw"/>
        <xsl:text>"noop":{"wt":"noop"}</xsl:text>
        <xsl:text>},"i":0}}]}</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation[@publication-type='book'] | mixed-citation[@publication-type='book']" mode="data-mw">
        <!-- {"parts":[{"template":{"target":{"wt":"cite book","href":"./Template:Cite_book"},"params":{"first":{"wt":"foo"},"last":{"wt":"bar"},"title":{"wt":"wibble"}},"i":0}}]} -->
        <xsl:text>{"parts":[{"template":{"target":{"wt":"cite book","href":"./Template:Cite_book"},"params":{</xsl:text>
        <xsl:apply-templates select="article-title | year | collab | source | volume | publisher-name | fpage" mode="data-mw"/>
        <xsl:text>"noop":{"wt":"noop"}</xsl:text>
        <xsl:text>},"i":0}}]}</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation/article-title | mixed-citation/article-title" mode="data-mw">
        <xsl:text>"title":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation/year | mixed-citation/year" mode="data-mw">
        <xsl:text>"date":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation/collab | mixed-citation/collab" mode="data-mw">
        <xsl:text>"publisher":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation[@publication-type='journal']/source | mixed-citation[@publication-type='journal']/source" mode="data-mw">
        <xsl:text>"journal":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation[@publication-type='journal']/volume | mixed-citation[@publication-type='journal']/volume" mode="data-mw">
        <xsl:text>"volume":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation/publisher-name | mixed-citation/publisher-name" mode="data-mw">
        <xsl:text>"publisher":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation[@publication-type='journal']/fpage | mixed-citation[@publication-type='journal']/fpage" mode="data-mw">
        <xsl:text>"page":{"wt":"</xsl:text>
        <xsl:value-of select="."/>
        <xsl:text>"},</xsl:text>
    </xsl:template>

    <xsl:template match="element-citation/* | mixed-citation/*" mode="data-mw">
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

  <!-- ============================================================= -->
  <!--  TABLES                                                       -->
  <!-- ============================================================= -->
  <!--  Tables are already in XHTML, and can simply be copied
        through.

        Technically parsoidHTML only deals with tbody when generating
        tables from wikitext, but it can consume them fine, so I'm
        going to leave them in when doing the conversion.            -->


  <xsl:template match="table | thead | tbody | tfoot | tr | th | td">
    <xsl:copy>
      <xsl:apply-templates select="@rowspan | @colspan" mode="table-copy"/>
      <xsl:apply-templates/>
    </xsl:copy>
  </xsl:template>

  <xsl:template match="array/tbody">
    <table>
      <xsl:copy>
      <xsl:apply-templates select="@*" mode="table-copy"/>
      <xsl:apply-templates/>
    </xsl:copy>
    </table>
  </xsl:template>

  <xsl:template match="@*" mode="table-copy">
    <xsl:copy-of select="."/>
  </xsl:template>

  <xsl:template match="@content-type" mode="table-copy"/>


</xsl:stylesheet>
