FROM scratch

MAINTAINER Christoph Görn <goern@b4mad.net>

LABEL Component="bn-o" \
      Name="goern/bn-o-0-generic" \
      Version="0.1.0" \
      Release="1"

LABEL io.k8s.description="This is bn-o!" \
      io.k8s.display-name="bn-o 0.1.0" \
      io.openshift.tags="bn-o"

# add bn-o to the container image
ADD bn-o /

# the entrypoint
ENTRYPOINT ["/bn-o"]
