##########
#
# This Script generates the MD Files for Nexus Data Sources and Resource Objects if not present
# and populates the nexus.erb layout file accourdingly
#
##########

#!/bin/bash

cat > website/nexus.erb << EOI
<% wrap_layout :inner do %>
  <% content_for :sidebar do %>
    <div class="docs-sidebar hidden-print affix-top" role="complementary">
      <ul class="nav docs-sidenav">
        <li<%= sidebar_current("docs-home") %>>
        <a href="/docs/providers/index.html">All Providers</a>
        </li>

        <li<%= sidebar_current("docs-nexus-index") %>>
        <a href="/docs/providers/nexus/index.html">Nexus Provider</a>
        </li>
        <li<%= sidebar_current("docs-nexus-datasource") %>>
          <a href="#">Data Sources</a>
          <ul class="nav nav-visible">
EOI

writingDataSources=true

for res in $(ls nexus/data_source_* nexus/resource_* | grep -v test | sort )
do 

  res=$(echo $res | awk -F/ '{print $2}' | sed 's/\.go//g')

  if [[ "$res" =~ "resource" ]]
  then
    typeShort="r"
    typeLong="resource"
    if [[ "${writingDataSources}" == true ]]
    then
      cat >> website/nexus.erb << EOI
          </ul>
        </li>
        <li<%= sidebar_current("docs-nexus-resource") %>>
        <a href="#">Resources</a>
        <ul class="nav nav-visible">
EOI
      writingDataSources=false
    fi
  else
    typeShort="d"
    typeLong="data"
  fi

  echo $res

  if [[ ! -f website/docs/$typeShort/$res.html.markdown ]]
  then
    cat > website/docs/$typeShort/${res}.html.markdown << EOI
---
layout: "nexus"
page_title: "Nexus: $res"
sidebar_current: "docs-nexus-$typeLong-source"
description: |-
  Sample $typeLong source in the Terraform provider Nexus.
---

# $res

Sample data source in the Terraform provider scaffolding.

## Example Usage

\`\`\`hcl
data "scaffolding_data_source" "example" {
  sample_attribute = "foo"
}
\`\`\`

## Attributes Reference

* sample_attribute - Sample attribute.

EOI
  fi

  cat >> website/nexus.erb << EOI
            <li<%= sidebar_current("$res") %>>
              <a href="/docs/providers/nexus/$typeShort/${res}.html">$res</a>
            </li>
EOI
done

cat >> website/nexus.erb << EOI
        </ul>
        </li>
      </ul>
    </div>
  <% end %>

  <%= yield %>
  <% end %>
EOI
